package collection

import (
	"time"

	"github.com/scjalliance/drivestream/page"
)

// A Writer writes data for a collection to a repository.
type Writer struct {
	repo      Repository
	seqNum    SeqNum
	nextState StateNum
	nextPage  page.SeqNum
	instance  string
}

// NewWriter returns a collection writer for the given sequence number.
func NewWriter(repo Repository, seqNum SeqNum, instance string) (*Writer, error) {
	nextState, err := repo.NextCollectionState(seqNum)
	if err != nil {
		return nil, err
	}

	nextPage, err := repo.NextPage(seqNum)
	if err != nil {
		return nil, err
	}

	return &Writer{
		repo:      repo,
		seqNum:    seqNum,
		nextState: nextState,
		nextPage:  nextPage,
		instance:  instance,
	}, nil
}

// Data returns information about the collection.
func (w *Writer) Data() (Data, error) {
	var buf [1]Data
	_, err := w.repo.Collections(w.seqNum, buf[:])
	return buf[0], err
}

// NextState returns the state number of the next state to be written.
func (w *Writer) NextState() StateNum {
	return w.nextState
}

// LastState returns the last state of the collection.
func (w *Writer) LastState() (State, error) {
	var buf [1]State
	_, err := w.repo.CollectionStates(w.seqNum, w.nextState-1, buf[:])
	return buf[0], err
}

// SetState sets the state of the collection.
func (w *Writer) SetState(phase Phase, pageNum page.SeqNum) error {
	err := w.repo.CreateCollectionState(w.seqNum, w.nextState, State{
		Time:     time.Now(),
		Instance: w.instance,
		StateData: StateData{
			Phase: phase,
			Page:  pageNum,
		},
	})
	if err == nil {
		w.nextState++
	}
	return err
}

// NextPage returns the page number of the next page to be written
func (w *Writer) NextPage() page.SeqNum {
	return w.nextPage
}

// LastPage reads the last page from the collection.
func (w *Writer) LastPage() (page.Data, error) {
	var buf [1]page.Data
	_, err := w.repo.Pages(w.seqNum, w.nextPage-1, buf[:])
	return buf[0], err
}

// AddPage adds the page to the collection.
func (w *Writer) AddPage(data page.Data) error {
	err := w.repo.CreatePage(w.seqNum, w.nextPage, data)
	if err == nil {
		w.nextPage++
	}
	return err
}

// ClearPages removes pages affiliated with the collection.
func (w *Writer) ClearPages() error {
	err := w.repo.ClearPages(w.seqNum)
	if err == nil {
		w.nextPage = 0
	}
	return err
}
