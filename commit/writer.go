package commit

import (
	"time"

	"github.com/scjalliance/drivestream/page"
)

// A Writer writes data for a commit to a repository.
type Writer struct {
	ref       Reference
	nextState StateNum
	instance  string
	time      time.Time
}

// NewWriter returns a commit writer for the given commit.
func NewWriter(ref Reference, instance string) (*Writer, error) {
	nextState, err := ref.States().Next()
	if err != nil {
		return nil, err
	}

	return &Writer{
		ref:       ref,
		nextState: nextState,
		instance:  instance,
	}, nil
}

// Data returns information about the commit.
func (w *Writer) Data() (Data, error) {
	return w.reader().Data()
}

// NextState returns the state number of the next state to be written.
func (w *Writer) NextState() StateNum {
	return w.reader().NextState()
}

// LastState returns the last state of the commit.
func (w *Writer) LastState() (State, error) {
	return w.reader().LastState()
}

// State returns the requested state from the commit.
func (w *Writer) State(stateNum StateNum) (State, error) {
	return w.reader().State(stateNum)
}

// States returns a slice of all states of the commit in ascending
// order.
func (w *Writer) States() ([]State, error) {
	return w.reader().States()
}

// SetState sets the state of the commit.
func (w *Writer) SetState(phase Phase, pageNum page.SeqNum) error {
	err := w.ref.State(w.nextState).Create(State{
		Time:     time.Now().UTC(),
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

// reader returns a Reader for w.
func (w *Writer) reader() *Reader {
	return &Reader{
		ref:       w.ref,
		nextState: w.nextState,
	}
}
