package collection

import "github.com/scjalliance/drivestream/page"

// Reader provides readonly access to a collection.
type Reader struct {
	ref       Reference
	nextState StateNum
	nextPage  page.SeqNum
}

// NewReader returns a collection reader for the given sequence number.
func NewReader(ref Reference) (*Reader, error) {
	nextState, err := ref.States().Next()
	if err != nil {
		return nil, err
	}

	nextPage, err := ref.Pages().Next()
	if err != nil {
		return nil, err
	}

	return &Reader{
		ref:       ref,
		nextState: nextState,
		nextPage:  nextPage,
	}, nil
}

// Data returns information about the collection.
func (r *Reader) Data() (Data, error) {
	return r.ref.Data()
}

// NextState returns the state number of the next state to be written.
func (r *Reader) NextState() StateNum {
	return r.nextState
}

// LastState returns the last state of the collection.
func (r *Reader) LastState() (State, error) {
	return r.State(r.nextState - 1)
}

// State returns the requested state from the collection.
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
		return nil, StatesTruncated{Drive: r.ref.Drive(), Collection: r.ref.SeqNum()}
	}
	return states, err
}

// NextPage returns the page number of the next page to be written
func (r *Reader) NextPage() page.SeqNum {
	return r.nextPage
}

// LastPage returns the last page from the collection.
func (r *Reader) LastPage() (page.Data, error) {
	return r.Page(r.nextPage - 1)
}

// Page returns the requested page from the collection.
func (r *Reader) Page(pageNum page.SeqNum) (page.Data, error) {
	return r.ref.Page(pageNum).Data()
}

// Pages returns a slice of all pages within the collection in ascending
// order.
//
// Note that this may allocate a significant amount of memory for large
// collections.
//
// TODO: Consider making this a buffered call.
func (r *Reader) Pages() ([]page.Data, error) {
	if r.nextPage == 0 {
		return nil, nil
	}
	pages := make([]page.Data, r.nextPage)
	n, err := r.ref.Pages().Read(0, pages)
	if err != nil {
		return nil, err
	}
	if n != len(pages) {
		return nil, PagesTruncated{Drive: r.ref.Drive(), Collection: r.ref.SeqNum()}
	}
	return pages, err
}
