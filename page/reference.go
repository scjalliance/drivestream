package page

// Reference provides access to a collection state.
type Reference interface {
	// SeqNum returns the sequence number of the page.
	SeqNum() SeqNum

	// Create creates the page with the given data.
	Create(data Data) error

	// Data returns the page data.
	Data() (Data, error)
}
