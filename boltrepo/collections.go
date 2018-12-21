package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.Sequence = (*Collections)(nil)

// Collections accesses a sequence of collections in a bolt repository.
type Collections struct {
	db    *bolt.DB
	drive resource.ID
}

// Path returns the path of the collections.
func (ref Collections) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, ref.drive.String(), CollectionBucket}
}

// Next returns the sequence number to use for the next collection.
func (ref Collections) Next() (n collection.SeqNum, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		collections := collectionsBucket(tx, ref.drive)
		if collections == nil {
			return nil
		}
		cursor := collections.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return nil
		}
		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadCollectionKey{Drive: ref.drive, BadKey: key}
		}
		n = collection.SeqNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads collection data for a range of collections
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (ref Collections) Read(start collection.SeqNum, p []collection.Data) (n int, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		collections := collectionsBucket(tx, ref.drive)
		if collections == nil {
			return collection.NotFound{Drive: ref.drive, Collection: start}
		}
		cursor := collections.Cursor()
		pos := start
		key := makeCollectionKey(pos)
		k, _ := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.NotFound{Drive: ref.drive, Collection: start}
		}
		for n < len(p) {
			col := collections.Bucket(k)
			if col == nil {
				return collection.DataInvalid{Drive: ref.drive, Collection: pos} // All collections must be buckets
			}
			value := col.Get([]byte(DataKey))
			if value == nil {
				return collection.DataInvalid{Drive: ref.drive, Collection: pos}
			}
			if err := json.Unmarshal(value, &p[n]); err != nil {
				// TODO: Wrap the error in DataInvalid?
				return err
			}
			n++
			k, _ = cursor.Next()
			if k == nil {
				break
			}
			pos = start + collection.SeqNum(n)
			key = makeCollectionKey(pos)
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

// Ref returns a collection reference.
func (ref Collections) Ref(c collection.SeqNum) collection.Reference {
	return Collection{
		db:         ref.db,
		drive:      ref.drive,
		collection: c,
	}
}
