package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ page.Sequence = (*Pages)(nil)

// Pages accesses a sequence of pages in a bolt repository.
type Pages struct {
	db         *bolt.DB
	drive      resource.ID
	collection collection.SeqNum
}

// Path returns the path of the pages.
func (ref Pages) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, ref.drive.String(), CollectionBucket, ref.collection.String(), PageBucket}
}

// Next returns the sequence number to use for the next page of the
// collection.
func (ref Pages) Next() (n page.SeqNum, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		pages := col.Bucket([]byte(PageBucket))
		if pages == nil {
			return nil
		}
		cursor := pages.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return nil
		}
		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadPageKey{Drive: ref.drive, Collection: ref.collection, BadKey: key}
		}
		n = page.SeqNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads the requested pages from a collection.
func (ref Pages) Read(start page.SeqNum, p []page.Data) (n int, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		pages := col.Bucket([]byte(PageBucket))
		if pages == nil {
			return collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: start}
		}
		cursor := pages.Cursor()
		pos := start
		key := makePageKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: start}
		}
		for n < len(p) {
			if v == nil {
				return collection.PageDataInvalid{Drive: ref.drive, Collection: ref.collection, Page: pos} // All pages must be non-nil
			}
			if err := json.Unmarshal(v, &p[n]); err != nil {
				// TODO: Wrap the error in PageDataInvalid?
				return err
			}
			n++
			k, v = cursor.Next()
			if k == nil {
				break
			}
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadPageKey{Drive: ref.drive, Collection: ref.collection, BadKey: key}
			}
			pos = start + page.SeqNum(n)
			key = makePageKey(pos)
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

// Ref returns a page reference for the sequence number.
func (ref Pages) Ref(seqNum page.SeqNum) page.Reference {
	return Page{
		db:         ref.db,
		drive:      ref.drive,
		collection: ref.collection,
		page:       seqNum,
	}
}

// Clear removes all pages affiliated with a collection.
func (ref Pages) Clear() error {
	return ref.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		key := []byte(PageBucket)
		pages := col.Bucket(key)
		if pages == nil {
			return nil
		}
		return col.DeleteBucket(key)
	})
}
