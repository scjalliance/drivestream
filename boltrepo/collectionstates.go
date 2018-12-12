package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.StateSequence = (*CollectionStates)(nil)

// CollectionStates accesses a sequence of collection states in a
// bolt repository.
type CollectionStates struct {
	db         *bolt.DB
	drive      resource.ID
	collection collection.SeqNum
}

// Next returns the state number to use for the next state.
func (seq CollectionStates) Next() (n collection.StateNum, err error) {
	err = seq.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, seq.drive, seq.collection)
		if col == nil {
			return collection.NotFound{Drive: seq.drive, Collection: seq.collection}
		}
		states := col.Bucket([]byte(StateBucket))
		if states == nil {
			return nil
		}
		cursor := states.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return nil
		}
		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadCollectionStateKey{Drive: seq.drive, Collection: seq.collection, BadKey: key}
		}
		n = collection.StateNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads a subset of states from the sequence, starting at start.
// Up to len(p) states will be returned in p. The number of states
// returned is provided as n.
func (seq CollectionStates) Read(start collection.StateNum, p []collection.State) (n int, err error) {
	err = seq.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, seq.drive, seq.collection)
		if col == nil {
			return collection.NotFound{Drive: seq.drive, Collection: seq.collection}
		}
		states := col.Bucket([]byte(StateBucket))
		if states == nil {
			return collection.StateNotFound{Drive: seq.drive, Collection: seq.collection, State: start}
		}
		cursor := states.Cursor()
		pos := start
		key := makeCollectionStateKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.StateNotFound{Drive: seq.drive, Collection: seq.collection, State: start}
		}
		for n < len(p) {
			if v == nil {
				return collection.StateInvalid{Drive: seq.drive, Collection: seq.collection, State: pos} // All states must be non-nil
			}
			if err := json.Unmarshal(v, &p[n]); err != nil {
				// TODO: Wrap the error in an InvalidState?
				return err
			}
			n++
			k, v = cursor.Next()
			if k == nil {
				break
			}
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCollectionStateKey{Drive: seq.drive, Collection: seq.collection, BadKey: key}
			}
			pos = start + collection.StateNum(n)
			key = makeCollectionStateKey(pos)
			if !bytes.Equal(key[:], k) {
				// The next key doesn't match the expected sequence number
				// TODO: Consider returning an error here?
				break
			}
		}
		return nil
	})
	return n, err
}

// Ref returns a collection state reference for the sequence number.
func (seq CollectionStates) Ref(stateNum collection.StateNum) collection.StateReference {
	return CollectionState{
		db:         seq.db,
		drive:      seq.drive,
		collection: seq.collection,
		state:      stateNum,
	}
}
