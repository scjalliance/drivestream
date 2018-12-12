package commit

// A StateSequence is an ordered series of drivestream commit states.
type StateSequence interface {
	// Next returns the sequence number to use for the next state.
	Next() (n StateNum, err error)

	// Read reads a subset of states from the sequence, starting at start.
	// Up to len(p) states will be returned in p. The number of states
	// returned is provided as n.
	Read(start StateNum, p []State) (n int, err error)

	// Ref returns a commit state reference for the sequence number.
	Ref(stateNum StateNum) StateReference
}
