package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ page.Reference = (*Page)(nil)

// Page is a drivestream page reference for a bolt repository.
type Page struct {
	db         *bolt.DB
	drive      resource.ID
	collection collection.SeqNum
	page       page.SeqNum
}

// SeqNum returns the sequence number of the page.
func (ref Page) SeqNum() page.SeqNum {
	return ref.page
}

// Create creates the page with the given data.
func (ref Page) Create(data page.Data) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		pages, err := col.CreateBucketIfNotExists([]byte(PageBucket))
		if err != nil {
			return err
		}

		var expected page.SeqNum
		{
			cursor := pages.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadPageKey{Drive: ref.drive, Collection: ref.collection, BadKey: key}
			default:
				expected = page.SeqNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.page != expected {
			return collection.PageOutOfOrder{Drive: ref.drive, Collection: ref.collection, Page: ref.page, Expected: expected}
		}

		key := makePageKey(ref.page)
		return pages.Put(key[:], value)
	})
}

// Data returns the page data.
func (ref Page) Data() (data page.Data, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, ref.drive, ref.collection)
		if col == nil {
			return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
		}
		pages := col.Bucket([]byte(PageBucket))
		if pages == nil {
			return collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: ref.page}
		}
		key := makePageKey(ref.page)
		value := pages.Get(key[:])
		if value == nil {
			return collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: ref.page}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}
