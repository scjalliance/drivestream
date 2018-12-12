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
func (seq CollectionStates) Next() (n collection.StateNum, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, nil
	}
	if seq.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, nil
	}
	return collection.StateNum(len(drv.Collections[seq.collection].States)), nil
}

// Read reads a subset of states from the sequence, starting at start.
// Up to len(p) states will be returned in p. The number of states
// returned is provided as n.
func (seq CollectionStates) Read(start collection.StateNum, p []collection.State) (n int, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	if seq.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	length := collection.StateNum(len(drv.Collections[seq.collection].States))
	if start >= length {
		return 0, collection.StateNotFound{Drive: seq.drive, Collection: seq.collection, State: start}
	}
	for n < len(p) && start+collection.StateNum(n) < length {
		p[n] = drv.Collections[seq.collection].States[start+collection.StateNum(n)]
		n++
	}
	return n, nil
}

// Ref returns a collection state reference for the sequence number.
func (seq CollectionStates) Ref(stateNum collection.StateNum) collection.StateReference {
	return CollectionState{
		repo:       seq.repo,
		drive:      seq.drive,
		collection: seq.collection,
		state:      stateNum,
	}
}
