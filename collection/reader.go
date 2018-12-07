package collection

import "github.com/scjalliance/drivestream/page"

// Reader provides readonly access to a collection.
type Reader struct {
	repo      Repository
	seqNum    SeqNum
	nextState StateNum
	nextPage  page.SeqNum
}

// NewReader returns a collection reader for the given sequence number.
func NewReader(repo Repository, seqNum SeqNum) (*Reader, error) {
	nextState, err := repo.NextCollectionState(seqNum)
	if err != nil {
		return nil, err
	}

	nextPage, err := repo.NextPage(seqNum)
	if err != nil {
		return nil, err
	}

	return &Reader{
		repo:      repo,
		seqNum:    seqNum,
		nextState: nextState,
		nextPage:  nextPage,
	}, nil
}

// Data returns information about the collection.
func (r *Reader) Data() (Data, error) {
	var buf [1]Data
	_, err := r.repo.Collections(r.seqNum, buf[:])
	return buf[0], err
}

// NextState returns the state number of the next state to be written.
func (r *Reader) NextState() StateNum {
	return r.nextState
}

// LastState returns the last state of the collection.
func (r *Reader) LastState() (State, error) {
	var buf [1]State
	_, err := r.repo.CollectionStates(r.seqNum, r.nextState-1, buf[:])
	return buf[0], err
}

// States returns a slice of all states of the collection in ascending
// order.
func (r *Reader) States() ([]State, error) {
	if r.nextState == 0 {
		return nil, nil
	}
	states := make([]State, r.nextState)
	n, err := r.repo.CollectionStates(r.seqNum, 0, states)
	if err != nil {
		return nil, err
	}
	if n != len(states) {
		return nil, TruncatedStates{SeqNum: 0}
	}
	return states, err
}

// NextPage returns the page number of the next page to be written
func (r *Reader) NextPage() page.SeqNum {
	return r.nextPage
}

// LastPage reads the last page from the collection.
func (r *Reader) LastPage() (page.Data, error) {
	var buf [1]page.Data
	_, err := r.repo.Pages(r.seqNum, r.nextPage-1, buf[:])
	return buf[0], err
}

// Pages returns a slice of all pages within the collection in ascending
// order.
//
// Note that this may allocate a significant amount of memory for large
// collections.
//
// TODO: Consider making this a buffered call.
func (r *Reader) Pages() ([]page.Data, error) {
	if r.nextState == 0 {
		return nil, nil
	}
	pages := make([]page.Data, r.nextPage)
	n, err := r.repo.Pages(r.seqNum, 0, pages)
	if err != nil {
		return nil, err
	}
	if n != len(pages) {
		return nil, TruncatedPages{SeqNum: 0}
	}
	return pages, err
}
