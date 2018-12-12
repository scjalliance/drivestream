package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.Reference = (*Collection)(nil)

// Collection is a drivestream collection reference for an in-memory
// repository.
type Collection struct {
	repo       *Repository
	drive      resource.ID
	collection collection.SeqNum
}

// Drive returns the drive ID of the collection.
func (ref Collection) Drive() resource.ID {
	return ref.drive
}

// SeqNum returns the sequence number of the collection.
func (ref Collection) SeqNum() collection.SeqNum {
	return ref.collection
}

// Exists returns true if the collection exists.
func (ref Collection) Exists() (exists bool, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return false, nil
	}
	return ref.collection < collection.SeqNum(len(drv.Collections)), nil
}

// Create creates a new collection with the given sequence number and data.
// If a collection already exists with the sequence number an error will be
// returned.
func (ref Collection) Create(data collection.Data) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		drv = newDriveEntry()
	}
	expected := collection.SeqNum(len(drv.Collections))
	if ref.collection != expected {
		return collection.OutOfOrder{Drive: ref.drive, Collection: ref.collection, Expected: expected}
	}
	drv.Collections = append(drv.Collections, newCollectionEntry(data))
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns information about the collection.
func (ref Collection) Data() (data collection.Data, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return collection.Data{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.Data{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	return drv.Collections[ref.collection].Data, nil
}

// States returns the state sequence for the collection.
func (ref Collection) States() collection.StateSequence {
	return CollectionStates{
		repo:       ref.repo,
		drive:      ref.drive,
		collection: ref.collection,
	}
}

// State returns a state reference.
func (ref Collection) State(stateNum collection.StateNum) collection.StateReference {
	return ref.States().Ref(stateNum)
}

// Pages returns the page sequence for the collection.
func (ref Collection) Pages() page.Sequence {
	return Pages{
		repo:       ref.repo,
		drive:      ref.drive,
		collection: ref.collection,
	}
}

// Page returns a page reference.
func (ref Collection) Page(pageNum page.SeqNum) page.Reference {
	return ref.Pages().Ref(pageNum)
}
