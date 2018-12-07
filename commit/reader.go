package commit

// Reader provides readonly access to a commit.
type Reader struct {
	repo      Repository
	seqNum    SeqNum
	nextState StateNum
}

// NewReader returns a commit reader for the given sequence number.
func NewReader(repo Repository, seqNum SeqNum) (*Reader, error) {
	nextState, err := repo.NextCommitState(seqNum)
	if err != nil {
		return nil, err
	}

	return &Reader{
		repo:      repo,
		seqNum:    seqNum,
		nextState: nextState,
	}, nil
}

// Data returns information about the commit.
func (r *Reader) Data() (Data, error) {
	var buf [1]Data
	_, err := r.repo.Commits(r.seqNum, buf[:])
	return buf[0], err
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
	var buf [1]State
	_, err := r.repo.CommitStates(r.seqNum, stateNum, buf[:])
	return buf[0], err
}

// States returns a slice of all states of the commit in ascending
// order.
func (r *Reader) States() ([]State, error) {
	if r.nextState == 0 {
		return nil, nil
	}
	states := make([]State, r.nextState)
	n, err := r.repo.CommitStates(r.seqNum, 0, states)
	if err != nil {
		return nil, err
	}
	if n != len(states) {
		return nil, TruncatedStates{SeqNum: 0}
	}
	return states, err
}
