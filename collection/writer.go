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
	return w.reader().Data()
}

// NextState returns the state number of the next state to be written.
func (w *Writer) NextState() StateNum {
	return w.reader().NextState()
}

// LastState returns the last state of the collection.
func (w *Writer) LastState() (State, error) {
	return w.reader().LastState()
}

// State returns the requested state from the collection.
func (w *Writer) State(stateNum StateNum) (State, error) {
	return w.reader().State(stateNum)
}

// States returns a slice of all states of the collection in ascending
// order.
func (w *Writer) States() ([]State, error) {
	return w.reader().States()
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
	return w.reader().NextPage()
}

// LastPage reads the last page from the collection.
func (w *Writer) LastPage() (page.Data, error) {
	return w.reader().LastPage()
}

// Page returns the requested page from the collection.
func (w *Writer) Page(pageNum page.SeqNum) (page.Data, error) {
	return w.reader().Page(pageNum)
}

// Pages returns a slice of all pages within the collection in ascending
// order.
//
// Note that this may allocate a significant amount of memory for large
// collections.
//
// TODO: Consider making this a buffered call.
func (w *Writer) Pages() ([]page.Data, error) {
	return w.reader().Pages()
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

// reader returns a Reader for w.
func (w *Writer) reader() *Reader {
	return &Reader{
		repo:      w.repo,
		seqNum:    w.seqNum,
		nextState: w.nextState,
		nextPage:  w.nextPage,
	}
}
