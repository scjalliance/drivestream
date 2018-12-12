package collection

// A StateSequence is an ordered series of drivestream collection states.
type StateSequence interface {
	// Next returns the state number to use for the next state.
	Next() (n StateNum, err error)

	// Read reads a subset of states from the sequence, starting at start.
	// Up to len(p) states will be returned in p. The number of states
	// returned is provided as n.
	Read(start StateNum, p []State) (n int, err error)

	// Ref returns a collection state reference for the state number.
	Ref(stateNum StateNum) StateReference
}
