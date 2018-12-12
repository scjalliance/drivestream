package commit

// StateReference is a reference to a commit state.
type StateReference interface {
	// StateNum returns the state number of the reference.
	StateNum() StateNum

	// Create creates the commit state with the given data. If a state already
	// exists with the state number an error will be returned.
	Create(data State) error

	// Data returns the commit state data.
	Data() (State, error)
}
