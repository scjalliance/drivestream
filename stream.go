package drivestream

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

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
	update := newTaskLogger(s.stdout).Task(fmt.Sprintf("DRIVE %s", s.repo.DriveID()))

	update.Log("Update started\n")
	defer func(e *error) {
		if *e != nil {
			update.Log("ERROR: %v\n", *e)
			update.Log("Update aborted after %s\n", update.Duration())
		} else {
			update.Log("Update finished in %s\n", update.Duration())
		}
	}(&err)

	update.Log("Retrieving existing collections from the repository\n")
	seqNum, err := s.repo.NextCollection()
	if err != nil {
		return err
	}

	switch {
	case seqNum > 0:
		if seqNum == 1 {
			update.Log("%d collection found\n", seqNum)
		} else {
			update.Log("%d collections found\n", seqNum)
		}
		seqNum--
	case seqNum < 0:
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
		preparation := col.Task("PREPARATION")

		preparation.Log("Creating writer\n")

		w, err := collection.NewWriter(s.repo, seqNum, s.instance)
		if err != nil {
			return err
		}

		preparation.Log("Reading collection data\n")
		data, err := w.Data()
		if err != nil {
			return err
		}

		if w.NextState() == 0 {
			preparation.Log("Adding the initial collection state to the repository\n")
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
			return err
		}

		switch state.Phase {
		case collection.PhaseCommitProcessing:
			preparation.Log("[%s] %s phase\n", data.Type, state.Phase)
		default:
			preparation.Log("[%s] %s phase, page %d\n", data.Type, state.Phase, state.Page)
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
			timestamp := time.Now()
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

			phase.Log("Updating collection state\n")
			if err := w.SetState(collection.PhaseFileCollection, 0); err != nil {
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
				return nil
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
					timestamp = time.Now()
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

			phase.Log("Updating collection state\n")
			if err := w.SetState(collection.PhaseChangeCollection, 0); err != nil {
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
					return nil
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
					timestamp = time.Now()
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

			phase.Log("Updating collection state\n")
			if err := w.SetState(collection.PhaseCommitProcessing, 0); err != nil {
				return err
			}

			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case collection.PhaseCommitProcessing:
			phase := col.Task(strings.ToUpper(collection.PhaseCommitProcessing.String()))
			phase.Log("Starting phase\n")
			phase.Log("Updating collection state\n")
			if err := w.SetState(collection.PhaseFinalized, 0); err != nil {
				return err
			}
			phase.Log("Finished phase in %s\n", phase.Duration())

			fallthrough
		case collection.PhaseFinalized:
			nextSeqNum := seqNum + 1
			col := update.Task(fmt.Sprintf("COLLECTION %d", nextSeqNum))

			eval := col.Task("EVAL")
			eval.Log("Checking for new changes\n")

			var startToken string

			if w.NextPage() == 0 {
				startToken = data.StartToken
			} else {
				last, err := w.LastPage()
				if err != nil {
					return nil
				}

				switch last.Type {
				case page.DriveList, page.FileList:
					startToken = data.StartToken
				case page.ChangeList:
					startToken = last.NextStartToken
				default:
					return fmt.Errorf("the previous collection page was of unrecognized type %d", data.Type)
				}
			}

			if startToken == "" {
				return fmt.Errorf("failed to determine starting token of next collection")
			}

			var (
				n         int
				nextToken string
				buf       [1]resource.Change
			)
			n, nextToken, _, err = c.Changes(ctx, startToken, buf[:])
			if err != nil {
				return err
			}
			if n == 0 && nextToken != "" {
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

// Cursor returns a new cursor for s.
/*
func (s *Stream) Cursor() *Cursor {
	return &Cursor{
		stream: s,
	}
}
*/
