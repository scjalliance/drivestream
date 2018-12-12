package commit

// Reader provides readonly access to a commit.
type Reader struct {
	ref       Reference
	nextState StateNum
}

// NewReader returns a commit reader for the given sequence number.
func NewReader(ref Reference) (*Reader, error) {
	nextState, err := ref.States().Next()
	if err != nil {
		return nil, err
	}

	return &Reader{
		ref:       ref,
		nextState: nextState,
	}, nil
}

// Data returns information about the commit.
func (r *Reader) Data() (Data, error) {
	return r.ref.Data()
}

// NextState returns the state number of the next state to be written.
func (r *Reader) NextState() StateNum {
	return r.nextState
}

// LastState returns the last state of the commit.
func (r *Reader) LastState() (State, error) {
	return r.State(r.nextState - 1)
}

// State returns the requested state from the commit.
func (r *Reader) State(stateNum StateNum) (State, error) {
	return r.ref.State(stateNum).Data()
}

// States returns a slice of all states of the collection in ascending
// order.
func (r *Reader) States() ([]State, error) {
	if r.nextState == 0 {
		return nil, nil
	}
	states := make([]State, r.nextState)
	n, err := r.ref.States().Read(0, states)
	if err != nil {
		return nil, err
	}
	if n != len(states) {
		return nil, StatesTruncated{Drive: r.ref.Drive(), Commit: r.ref.SeqNum()}
	}
	return states, err
}
