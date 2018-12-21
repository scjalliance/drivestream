package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.StateReference = (*CollectionState)(nil)

// CollectionState is a drivestream collection state accessor for a
// bolt repository.
type CollectionState struct {
	db         *bolt.DB
	drive      resource.ID
	collection collection.SeqNum
	state      collection.StateNum
}

// Path returns the path of the collection state.
func (ref CollectionState) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, ref.drive.String(), CollectionBucket, ref.collection.String(), StateBucket, ref.state.String()}
}

// StateNum returns the sequence number of the reference.
func (ref CollectionState) StateNum() collection.StateNum {
	return ref.state
}

// Create creates the collection state with the given data.
// If a state already exists with the state's sequence number an error
// will be returned.
func (ref CollectionState) Create(data collection.State) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		states, err := col.CreateBucketIfNotExists([]byte(StateBucket))
		if err != nil {
			return err
		}

		var expected collection.StateNum
		{
			cursor := states.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCollectionStateKey{Drive: ref.drive, Collection: ref.collection, BadKey: key}
			default:
				expected = collection.StateNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.state != expected {
			return collection.StateOutOfOrder{Drive: ref.drive, Collection: ref.collection, State: ref.state, Expected: expected}
		}

		key := makeCollectionStateKey(ref.state)
		return states.Put(key[:], value)
	})
}

// Data returns the collection state data.
func (ref CollectionState) Data() (data collection.State, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		states := col.Bucket([]byte(StateBucket))
		if states == nil {
			return collection.StateNotFound{Drive: ref.drive, Collection: ref.collection, State: ref.state}
		}
		key := makeCollectionStateKey(ref.state)
		value := states.Get(key[:])
		if value == nil {
			return collection.StateNotFound{Drive: ref.drive, Collection: ref.collection, State: ref.state}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}
