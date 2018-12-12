package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.StateReference = (*CollectionState)(nil)

// CollectionState is a reference to a collection state.
type CollectionState struct {
	repo       *Repository
	drive      resource.ID
	collection collection.SeqNum
	state      collection.StateNum
}

// StateNum returns the state number of the reference.
func (ref CollectionState) StateNum() collection.StateNum {
	return ref.state
}

// Create creates the collection state with the given data. If a state already
// exists with the state number an error will be returned.
func (ref CollectionState) Create(data collection.State) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	expected := collection.StateNum(len(drv.Collections[ref.collection].States))
	if ref.state != expected {
		return collection.StateOutOfOrder{Drive: ref.drive, Collection: ref.collection, State: ref.state, Expected: expected}
	}
	drv.Collections[ref.collection].States = append(drv.Collections[ref.collection].States, data)
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns the collection state data.
func (ref CollectionState) Data() (data collection.State, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return collection.State{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.State{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.state >= collection.StateNum(len(drv.Collections[ref.collection].States)) {
		return collection.State{}, collection.StateNotFound{Drive: ref.drive, Collection: ref.collection, State: ref.state}
	}
	return drv.Collections[ref.collection].States[ref.state], nil
}
