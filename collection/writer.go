package collection

import (
	"time"

	"github.com/scjalliance/drivestream/page"
)

// A Writer writes data for a collection to a repository.
type Writer struct {
	ref       Reference
	nextState StateNum
	nextPage  page.SeqNum
	instance  string
}

// NewWriter returns a collection writer for the collection.
func NewWriter(ref Reference, instance string) (*Writer, error) {
	nextState, err := ref.States().Next()
	if err != nil {
		return nil, err
	}

	nextPage, err := ref.Pages().Next()
	if err != nil {
		return nil, err
	}

	return &Writer{
		ref:       ref,
		nextState: nextState,
		nextPage:  nextPage,
		instance:  instance,
	}, nil
}

// Data returns information about the collection.
func (w *Writer) Data() (Data, error) {
	return w.ref.Data()
}

// NextState returns the state number of the next state to be written.
func (w *Writer) NextState() StateNum {
	return w.nextState
}

// LastState returns the last state of the collection.
func (w *Writer) LastState() (State, error) {
	return w.ref.State(w.nextState - 1).Data()
}

// State returns the requested state from the collection.
func (w *Writer) State(stateNum StateNum) (State, error) {
	return w.ref.State(stateNum).Data()
}

// States returns a slice of all states of the collection in ascending
// order.
/*
func (w *Writer) States() ([]State, error) {
	return w.reader().States()
}
*/

// SetState sets the state of the collection.
func (w *Writer) SetState(phase Phase, pageNum page.SeqNum) error {
	err := w.ref.State(w.nextState).Create(State{
		Time:     time.Now().UTC(),
		Instance: w.instance,
		Phase:    phase,
		Page:     pageNum,
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
	return w.ref.Page(w.nextPage - 1).Data()
}

// Page returns the requested page from the collection.
func (w *Writer) Page(pageNum page.SeqNum) (page.Data, error) {
	return w.ref.Page(pageNum).Data()
}

// Pages returns a slice of all pages within the collection in ascending
// order.
//
// Note that this may allocate a significant amount of memory for large
// collections.
//
// TODO: Consider making this a buffered call.
/*
func (w *Writer) Pages() ([]page.Data, error) {
	return w.reader().Pages()
}
*/

// AddPage adds the page to the collection.
func (w *Writer) AddPage(data page.Data) error {
	err := w.ref.Page(w.nextPage).Create(data)
	if err == nil {
		w.nextPage++
	}
	return err
}

// ClearPages removes pages affiliated with the collection.
func (w *Writer) ClearPages() error {
	err := w.ref.Pages().Clear()
	if err == nil {
		w.nextPage = 0
	}
	return err
}

// reader returns a Reader for w.
/*
func (w *Writer) reader() *Reader {
	return &Reader{
		repo:      w.repo,
		seqNum:    w.seqNum,
		nextState: w.nextState,
		nextPage:  w.nextPage,
	}
}
*/
