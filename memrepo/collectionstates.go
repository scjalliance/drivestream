package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.StateSequence = (*CollectionStates)(nil)

// CollectionStates accesses a sequence of collection states in an in-memory
// repository.
type CollectionStates struct {
	repo       *Repository
	drive      resource.ID
	collection collection.SeqNum
}

// Next returns the state number to use for the next state.
func (ref CollectionStates) Next() (n collection.StateNum, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, nil
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, nil
	}
	return collection.StateNum(len(drv.Collections[ref.collection].States)), nil
}

// Read reads a subset of states from the sequence, starting at start.
// Up to len(p) states will be returned in p. The number of states
// returned is provided as n.
func (ref CollectionStates) Read(start collection.StateNum, p []collection.State) (n int, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	length := collection.StateNum(len(drv.Collections[ref.collection].States))
	if start >= length {
		return 0, collection.StateNotFound{Drive: ref.drive, Collection: ref.collection, State: start}
	}
	for n < len(p) && start+collection.StateNum(n) < length {
		p[n] = drv.Collections[ref.collection].States[start+collection.StateNum(n)]
		n++
	}
	return n, nil
}

// Ref returns a collection state reference for the sequence number.
func (ref CollectionStates) Ref(stateNum collection.StateNum) collection.StateReference {
	return CollectionState{
		repo:       ref.repo,
		drive:      ref.drive,
		collection: ref.collection,
		state:      stateNum,
	}
}
