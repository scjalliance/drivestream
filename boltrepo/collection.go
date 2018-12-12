package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.Reference = (*Collection)(nil)

// Collection is a drivestream collection reference for a bolt repository.
type Collection struct {
	db         *bolt.DB
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
	err = ref.db.View(func(tx *bolt.Tx) error {
		if collectionBucket(tx, ref.drive, ref.collection) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// Create creates a new collection with the given sequence number and data.
// If a collection already exists with the sequence number an error will be
// returned.
func (ref Collection) Create(data collection.Data) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		collections, err := createCollectionsBucket(tx, ref.drive)
		if err != nil {
			return err
		}

		var expected collection.SeqNum
		{
			cursor := collections.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCollectionKey{Drive: ref.drive, BadKey: key}
			default:
				expected = collection.SeqNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.collection != expected {
			return collection.OutOfOrder{Drive: ref.drive, Collection: ref.collection, Expected: expected}
		}

		key := makeCollectionKey(ref.collection)
		col, err := collections.CreateBucket(key[:])
		if err != nil {
			return err
		}
		return col.Put([]byte(DataKey), value)
	})
}

// Data returns information about the collection.
func (ref Collection) Data() (data collection.Data, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		value := col.Get([]byte(DataKey))
		if value == nil {
			return collection.DataInvalid{Drive: ref.drive, Collection: ref.collection}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}

// States returns the state sequence for the collection.
func (ref Collection) States() collection.StateSequence {
	return CollectionStates{
		db:         ref.db,
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
		db:         ref.db,
		drive:      ref.drive,
		collection: ref.collection,
	}
}

// Page returns a page reference.
func (ref Collection) Page(pageNum page.SeqNum) page.Reference {
	return ref.Pages().Ref(pageNum)
}
