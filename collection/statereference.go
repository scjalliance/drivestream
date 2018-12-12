package collection

// StateReference is a reference to a collection state.
type StateReference interface {
	// StateNum returns the sequence number of the reference.
	StateNum() StateNum

	// Create creates a new collection state with the given state number and data.
	// If a state already exists with the state number an error will be returned.
	Create(state State) error

	// Data returns the collection state.
	Data() (State, error)
}
