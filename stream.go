package drivestream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/scjalliance/drivestream/commit"

	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

const defaultPageSize = 1000

// Stream provides access to a stream of team drive changes.
type Stream struct {
	repo      Repository
	stdout    io.Writer
	collector Collector
	instance  string
	pageSize  int64
}

// New returns a new drive stream for the given service and team drive ID.
func New(repo Repository, options ...Option) *Stream {
	s := &Stream{
		repo:     repo,
		pageSize: defaultPageSize,
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// Update queries c for an updated set of changes, processes them and
// persists them in the stream's repository.
func (s *Stream) Update(ctx context.Context, c Collector) (err error) {
	update := newTaskLogger(s.stdout).Task(fmt.Sprintf("DRIVE %s", s.repo.DriveID())).Task("UPDATE")

	update.Log("Started  %s\n", time.Now().Format(time.RFC3339))
	defer func(e *error) {
		if *e != nil {
			update.Log("ERROR: %v\n", *e)
			update.Log("Aborted  %s | %s\n", time.Now().Format(time.RFC3339), update.Duration())
		} else {
			update.Log("Finished %s | %s\n", time.Now().Format(time.RFC3339), update.Duration())
		}
	}(&err)

	if err = s.collect(ctx, c, update); err != nil {
		return err
	}

	return s.buildCommits(ctx, update)
}

func (s *Stream) collect(ctx context.Context, c Collector, update taskLogger) (err error) {
	seqNum, err := s.repo.NextCollection()
	if err != nil {
		update.Log("Retrieving collections from the repository\n")
		return err
	}

	switch {
	case seqNum > 0:
		if seqNum == 1 {
			//update.Log("%d collection\n", seqNum)
		} else {
			//update.Log("%d collections\n", seqNum)
		}
		seqNum--
	case seqNum < 0:
		update.Log("Retrieving collections from the repository\n")
		return fmt.Errorf("the repository returned a negative number of collections (%d)", seqNum)
	case seqNum == 0:
		update.Log("No collections found. Initializing.\n")

		col := update.Task(fmt.Sprintf("COLLECTION %d", seqNum))
		init := col.Task("INIT")

		init.Log("Collecting starting change token\n")
		startToken, err := c.ChangeToken(ctx)
		if err != nil {
			return err
		}

		init.Log("Adding collection to the repository\n")
		data := collection.Data{
			Type:       collection.Full,
			StartToken: startToken,
		}
		if err = s.repo.CreateCollection(seqNum, data); err != nil {
			return err
		}
	}

	for {
		col := update.Task(fmt.Sprintf("COLLECTION %d", seqNum))
		eval := col.Task("EVAL")

		w, err := collection.NewWriter(s.repo, seqNum, s.instance)
		if err != nil {
			eval.Log("Creating writer\n")
			return err
		}

		data, err := w.Data()
		if err != nil {
			eval.Log("Reading collection data\n")
			return err
		}

		if w.NextState() == 0 {
			col.Task("INIT").Log("Adding the initial collection state to the repository\n")
			switch data.Type {
			case collection.Full:
				w.SetState(collection.PhaseDriveCollection, 0)
			case collection.Incremental:
				w.SetState(collection.PhaseChangeCollection, 0)
			default:
				return fmt.Errorf("unable to determine starting phase for unknown collection type %d", data.Type)
			}
		}

		state, err := w.LastState()
		if err != nil {
			eval.Log("Examining state\n")
			return err
		}

		if state.Page != 0 {
			eval.Log("%s | %s | PAGE %d\n", strings.ToUpper(data.Type.String()), strings.ToUpper(state.Phase.String()), state.Page)
		} else {
			eval.Log("%s | %s\n", strings.ToUpper(data.Type.String()), strings.ToUpper(state.Phase.String()))
		}

		switch state.Phase {
		case collection.PhaseDriveCollection:
			phase := col.Task(strings.ToUpper(collection.PhaseDriveCollection.String()))
			phase.Log("Starting phase\n")

			if w.NextPage() > 0 {
				phase.Log("Clearing previously written pages\n")
				if err = w.ClearPages(); err != nil {
					return err
				}
			}

			phase.Log("Collecting current drive data\n")
			timestamp := time.Now().UTC()
			record, err := c.Drive(ctx)
			if err != nil {
				return err
			}

			phase.Log("Adding drive data page %d to the repository\n", w.NextPage())
			pageData := page.Data{
				Type:      page.DriveList,
				Collected: timestamp,
				Changes:   []resource.Change{record},
			}
			if err := w.AddPage(pageData); err != nil {
				return err
			}

			if err := w.SetState(collection.PhaseFileCollection, 0); err != nil {
				phase.Log("Updating collection state\n")
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case collection.PhaseFileCollection:
			phase := col.Task(strings.ToUpper(collection.PhaseFileCollection.String()))
			phase.Log("Starting phase\n")

			if w.NextPage() == 0 {
				return errors.New("file collection pages must follow drive collection pages")
			}

			phase.Log("Retrieving the most recent page from the repository\n")
			last, err := w.LastPage()
			if err != nil {
				return err
			}

			var first bool
			var nextToken string
			switch last.Type {
			case page.DriveList:
				first = true
			case page.FileList:
				nextToken = last.NextPageToken
				if nextToken != "" {
					phase.Log("Resuming file data collection using saved token \"%s\"\n", nextToken)
				}
			case page.ChangeList:
				return errors.New("file collection pages must follow drive collection pages")
			default:
				return fmt.Errorf("the previous collection page was of unrecognized type %d", data.Type)
			}

			for first || nextToken != "" {
				first = false

				if nextToken == "" {
					phase.Log("Collecting file data\n")
				} else {
					phase.Log("Collecting file data with token: %s\n", nextToken)
				}

				var (
					n         int
					token     = nextToken
					buf       = make([]resource.Change, s.pageSize)
					timestamp = time.Now().UTC()
				)
				n, nextToken, err = c.Files(ctx, token, buf)
				if err != nil {
					return err
				}
				if n == 0 && nextToken != "" {
					return fmt.Errorf("the collector returned an empty file data page")
				}

				phase.Log("Adding file data page %d with %d entries to the repository\n", w.NextPage(), n)
				pageData := page.Data{
					Type:          page.FileList,
					Collected:     timestamp,
					PageToken:     token,
					NextPageToken: nextToken,
					Changes:       buf[:n],
				}
				if err := w.AddPage(pageData); err != nil {
					return err
				}
			}
			phase.Log("The end of the file data series has been reached\n")

			if err := w.SetState(collection.PhaseChangeCollection, 0); err != nil {
				phase.Log("Updating collection state\n")
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case collection.PhaseChangeCollection:
			phase := col.Task(strings.ToUpper(collection.PhaseChangeCollection.String()))
			phase.Log("Starting phase\n")

			var nextToken, nextStartToken string

			if w.NextPage() == 0 {
				nextToken = data.StartToken
			} else {
				phase.Log("Retrieving the most recent page from the repository\n")
				last, err := w.LastPage()
				if err != nil {
					return err
				}

				switch last.Type {
				case page.DriveList, page.FileList:
					nextToken = data.StartToken
				case page.ChangeList:
					nextToken = last.NextPageToken
					if nextToken != "" {
						phase.Log("Resuming change data collection using saved token \"%s\"\n", nextToken)
					}
				default:
					return fmt.Errorf("the previous collection page was of unrecognized type %d", data.Type)
				}
			}

			for nextToken != "" {
				phase.Log("Collecting change data with token: %s\n", nextToken)

				var (
					n         int
					token     = nextToken
					buf       = make([]resource.Change, s.pageSize)
					timestamp = time.Now().UTC()
				)
				n, nextToken, nextStartToken, err = c.Changes(ctx, token, buf[:])
				if err != nil {
					return err
				}
				if n == 0 && nextToken != "" {
					return fmt.Errorf("the collector returned an empty change data page")
				}

				if n > 0 || (nextStartToken != "" && nextStartToken != data.StartToken) {
					phase.Log("Adding change data page %d with %d entries to the repository\n", w.NextPage(), n)
					pageData := page.Data{
						Type:           page.ChangeList,
						Collected:      timestamp,
						PageToken:      token,
						NextPageToken:  nextToken,
						NextStartToken: nextStartToken,
						Changes:        buf[:n],
					}
					if err := w.AddPage(pageData); err != nil {
						return err
					}
				}
			}
			phase.Log("The end of the change data series has been reached\n")

			if err := w.SetState(collection.PhaseFinalized, 0); err != nil {
				phase.Log("Updating collection state\n")
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case collection.PhaseFinalized:
			nextSeqNum := seqNum + 1
			col := update.Task(fmt.Sprintf("COLLECTION %d", nextSeqNum))

			eval := col.Task("EVAL")

			var startToken string

			if w.NextPage() == 0 {
				startToken = data.StartToken
			} else {
				last, err := w.LastPage()
				if err != nil {
					eval.Log("Determining starting token\n")
					return err
				}

				switch last.Type {
				case page.DriveList, page.FileList:
					startToken = data.StartToken
				case page.ChangeList:
					startToken = last.NextStartToken
				default:
					eval.Log("Determining starting token\n")
					return fmt.Errorf("the previous collection page was of unrecognized type %d", data.Type)
				}
			}

			if startToken == "" {
				eval.Log("Determining starting token\n")
				return fmt.Errorf("failed to determine starting token of next collection")
			}

			var (
				n         int
				nextToken string
				buf       [1]resource.Change
			)
			n, nextToken, _, err = c.Changes(ctx, startToken, buf[:])
			if err != nil {
				eval.Log("Checking for new changes\n")
				return err
			}
			if n == 0 && nextToken != "" {
				eval.Log("Checking for new changes\n")
				return fmt.Errorf("the collector returned an empty change data page")
			}

			if n == 0 {
				eval.Log("No changes found\n")
				return nil
			}

			eval.Log("Changes found\n")

			init := col.Task("INIT")

			init.Log("Adding collection to the repository\n")
			data := collection.Data{
				Type:       collection.Incremental,
				StartToken: startToken,
			}
			if err = s.repo.CreateCollection(nextSeqNum, data); err != nil {
				return err
			}

			seqNum = nextSeqNum
		default:
			return fmt.Errorf("The collection phase is of unrecognized type %d", state.Phase)
		}
	}
}

func (s *Stream) buildCommits(ctx context.Context, update taskLogger) (err error) {
	seqNum, err := s.repo.NextCommit()
	if err != nil {
		update.Log("Retrieving commits from the repository\n")
		return err
	}

	switch {
	case seqNum > 0:
		if seqNum == 1 {
			//update.Log("%d commit\n", seqNum)
		} else {
			//update.Log("%d commits\n", seqNum)
		}
		seqNum--
	case seqNum < 0:
		update.Log("Retrieving commits from the repository\n")
		return fmt.Errorf("the repository returned a negative number of commits (%d)", seqNum)
	case seqNum == 0:
		isReady, err := s.readyToCommit(0)
		if err != nil {
			update.Log("Retrieving collections from the repository\n")
			return err
		}

		if !isReady {
			update.Log("Nothing to commit.\n")
			return nil
		}

		com := update.Task(fmt.Sprintf("COMMIT %d", seqNum))
		init := com.Task("INIT")
		init.Log("Adding commit to the repository\n")
		if err = s.repo.CreateCommit(seqNum, commit.Data{}); err != nil {
			return err
		}
	}

	for {
		com := update.Task(fmt.Sprintf("COMMIT %d", seqNum))
		eval := com.Task("EVAL")

		w, err := commit.NewWriter(s.repo, seqNum, s.instance)
		if err != nil {
			eval.Log("Creating writer\n")
			return err
		}

		data, err := w.Data()
		if err != nil {
			eval.Log("Reading commit data\n")
			return err
		}

		if w.NextState() == 0 {
			com.Task("INIT").Log("Adding commit state 0\n")
			w.SetState(commit.PhaseSourceProcessing, data.Source.Page)
		}

		state, err := w.LastState()
		if err != nil {
			eval.Log("Examining state\n")
			return err
		}

		var colType collection.Type
		{
			var buf [1]collection.Data
			if _, err := s.repo.Collections(data.Source.Collection, buf[:]); err != nil {
				eval.Log("Examining source collection\n")
				return err
			}
			colType = buf[0].Type
		}

		if state.Phase == commit.PhaseSourceProcessing {
			switch colType {
			case collection.Full:
				eval.Log("%s | COLLECTION %d [%s] | PAGE %d\n", strings.ToUpper(state.Phase.String()), data.Source.Collection, strings.ToUpper(colType.String()), state.Page)
			case collection.Incremental:
				eval.Log("%s | COLLECTION %d [%s] | PAGE %d | INDEX %d\n", strings.ToUpper(state.Phase.String()), data.Source.Collection, strings.ToUpper(colType.String()), data.Source.Page, data.Source.Index)
			}
		} else {
			eval.Log("%s\n", strings.ToUpper(state.Phase.String()))
		}

		switch state.Phase {
		case commit.PhaseSourceProcessing:
			phase := com.Task(strings.ToUpper(commit.PhaseSourceProcessing.String()))
			phase.Log("Starting phase\n")

			switch colType {
			case collection.Full:
				col, err := collection.NewReader(s.repo, data.Source.Collection)
				if err != nil {
					phase.Log("Examining source collection\n")
					return err
				}

				for {
					pg, err := col.Page(state.Page)
					if err != nil {
						phase.Log("Retrieving page data\n")
						return err
					}

					phase.Log("COL %d PAGE %d\n", data.Source.Collection, state.Page)

					for i := range pg.Changes {
						_ = pg.Changes[i]
						// TODO: Process changes
					}

					if state.Page+1 >= col.NextPage() {
						break
					}

					state.Page++
					if err := w.SetState(commit.PhaseSourceProcessing, state.Page); err != nil {
						phase.Log("Updating commit state\n")
						return err
					}
				}
			case collection.Incremental:
				phase.Log("COL %d PAGE %d INDEX %d\n", data.Source.Collection, data.Source.Page, data.Source.Index)
				// TODO: Process change
			default:
				return fmt.Errorf("the source collection is of unrecognized type %d", colType)
			}

			if err := w.SetState(commit.PhaseTreeProcessing, 0); err != nil {
				phase.Log("Updating commit state\n")
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case commit.PhaseTreeProcessing:
			phase := com.Task(strings.ToUpper(commit.PhaseTreeProcessing.String()))
			phase.Log("Starting phase\n")

			if err := w.SetState(commit.PhaseFinalized, 0); err != nil {
				phase.Log("Updating commit state\n")
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case commit.PhaseFinalized:
			nextSeqNum := seqNum + 1
			com := update.Task(fmt.Sprintf("COMMIT %d", nextSeqNum))

			eval := com.Task("EVAL")

			isReady := false
			nextSource := data.Source

			switch colType {
			case collection.Incremental:
				col, err := collection.NewReader(s.repo, nextSource.Collection)
				if err != nil {
					eval.Log("Examining source collection\n")
					return err
				}
				pg, err := col.Page(nextSource.Page)
				if err != nil {
					eval.Log("Retrieving page data\n")
					return err
				}
				switch {
				case nextSource.Index+1 < len(pg.Changes):
					isReady = true
					nextSource.Index++
				case nextSource.Page+1 < col.NextPage():
					isReady = true
					nextSource.Page++
					nextSource.Index = 0
				default:
					isReady, err = s.readyToCommit(nextSource.Collection + 1)
					if err != nil {
						eval.Log("Examining source collection\n")
						return err
					}
					nextSource.Collection++
					nextSource.Page = 0
					nextSource.Index = 0
				}
			case collection.Full:
				isReady, err = s.readyToCommit(nextSource.Collection + 1)
				if err != nil {
					eval.Log("Examining source collection\n")
					return err
				}
				nextSource.Collection++
				nextSource.Page = 0
				nextSource.Index = 0
			default:
				return fmt.Errorf("the source collection is of unrecognized type %d", colType)
			}

			if !isReady {
				eval.Log("No more data is ready for processing\n")
				return nil
			}

			eval.Log("More data is ready for processing\n")

			init := com.Task("INIT")

			init.Log("Adding commit to the repository\n")
			data := commit.Data{
				Source: nextSource,
			}
			if err = s.repo.CreateCommit(nextSeqNum, data); err != nil {
				return err
			}

			seqNum = nextSeqNum
		default:
			return fmt.Errorf("the commit phase is of unrecognized type %d", state.Phase)
		}
	}
}

func (s *Stream) readyToCommit(seqNum collection.SeqNum) (bool, error) {
	next, err := s.repo.NextCollection()
	if err != nil {
		return false, err
	}
	if seqNum >= next {
		return false, nil
	}
	col, err := collection.NewReader(s.repo, seqNum)
	if err != nil {
		return false, err
	}
	if col.NextState() == 0 {
		return false, nil
	}
	last, err := col.LastState()
	if err != nil {
		return false, err
	}
	return last.Phase == collection.PhaseFinalized, nil
}

// Cursor returns a new cursor for s.
/*
func (s *Stream) Cursor() *Cursor {
	return &Cursor{
		stream: s,
	}
}
*/
