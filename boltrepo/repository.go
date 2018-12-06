package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// Repository is a drive stream repository backed by a bolt database.
// It should be created by calling New.
type Repository struct {
	db          *bolt.DB
	teamDriveID resource.ID
}

// New returns a new drivestream bolt database for the team drive.
func New(db *bolt.DB, teamDriveID resource.ID) *Repository {
	return &Repository{
		db:          db,
		teamDriveID: teamDriveID,
	}
}

// DriveID returns the team drive ID of the repository.
func (repo *Repository) DriveID() resource.ID {
	return repo.teamDriveID
}

// NextCollection returns the sequence number to use for the next
// collection.
func (repo *Repository) NextCollection() (n collection.SeqNum, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		collections := collectionsBucket(tx, repo.teamDriveID)
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
			return BadCollectionKey{BadKey: key}
		}
		n = collection.SeqNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Collections returns collection data for a range of collections
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (repo *Repository) Collections(start collection.SeqNum, p []collection.Data) (n int, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		collections := collectionsBucket(tx, repo.teamDriveID)
		if collections == nil {
			return collection.NotFound{SeqNum: start}
		}
		cursor := collections.Cursor()
		pos := start
		key := makeCollectionKey(pos)
		k, _ := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.NotFound{SeqNum: start}
		}
		for n < len(p) {
			col := collections.Bucket(k)
			if col == nil {
				return collection.InvalidData{SeqNum: pos} // All collections must be buckets
			}
			value := col.Get([]byte(DataKey))
			if value == nil {
				return collection.InvalidData{SeqNum: pos}
			}
			if err := json.Unmarshal(value, &p[n]); err != nil {
				// TODO: Wrap the error in an InvalidData?
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

// CreateCollection creates a new collection with the given sequence
// number and data. If a collection already exists with the sequence
// number an error will be returned.
func (repo *Repository) CreateCollection(c collection.SeqNum, data collection.Data) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return repo.db.Update(func(tx *bolt.Tx) error {
		collections, err := createCollectionsBucket(tx, repo.teamDriveID)
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
				return BadCollectionKey{BadKey: key}
			default:
				expected = collection.SeqNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if c != expected {
			return collection.OutOfOrder{SeqNum: c, Expected: expected}
		}

		key := makeCollectionKey(c)
		col, err := collections.CreateBucket(key[:])
		if err != nil {
			return err
		}
		return col.Put([]byte(DataKey), value)
	})
}

// NextCollectionState returns the state number to use for the next
// state of the collection.
func (repo *Repository) NextCollectionState(c collection.SeqNum) (n collection.StateNum, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
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
			return BadStateKey{SeqNum: c, BadKey: key}
		}
		n = collection.StateNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// CollectionStates returns a range of collection states for the given
// collection, starting at the given state number. Up to len(p) states
// will be returned in p. The number of states is returned as n.
func (repo *Repository) CollectionStates(c collection.SeqNum, start collection.StateNum, p []collection.State) (n int, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
		}
		states := col.Bucket([]byte(StateBucket))
		if states == nil {
			return collection.StateNotFound{SeqNum: c, StateNum: start}
		}
		cursor := states.Cursor()
		pos := start
		key := makeStateKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.StateNotFound{SeqNum: c, StateNum: start}
		}
		for n < len(p) {
			if v == nil {
				return collection.InvalidState{SeqNum: c, StateNum: pos} // All collections must be non-nil
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
				return BadStateKey{SeqNum: c, BadKey: key}
			}
			pos = start + collection.StateNum(n)
			key = makeStateKey(pos)
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

// CreateCollectionState creates a new collection state with the given
// state number and data. If a state already exists with the state
// number an error will be returned.
func (repo *Repository) CreateCollectionState(c collection.SeqNum, stateNum collection.StateNum, state collection.State) error {
	value, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return repo.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
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
				return BadStateKey{SeqNum: c, BadKey: key}
			default:
				expected = collection.StateNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if stateNum != expected {
			return collection.StateOutOfOrder{SeqNum: c, StateNum: stateNum, Expected: expected}
		}

		key := makeStateKey(stateNum)
		return states.Put(key[:], value)
	})
}

// NextPage returns the sequence number to use for the next page of the
// collection.
func (repo *Repository) NextPage(c collection.SeqNum) (n page.SeqNum, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
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
			return BadPageKey{c, key}
		}
		n = page.SeqNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Pages returns the requested pages from a collection.
func (repo *Repository) Pages(c collection.SeqNum, start page.SeqNum, p []page.Data) (n int, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
		}
		pages := col.Bucket([]byte(PageBucket))
		if pages == nil {
			return collection.PageNotFound{SeqNum: c, PageNum: start}
		}
		cursor := pages.Cursor()
		pos := start
		key := makePageKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return collection.PageNotFound{SeqNum: c, PageNum: start}
		}
		for n < len(p) {
			if v == nil {
				return collection.InvalidPage{SeqNum: c, PageNum: pos} // All collections must be non-nil
			}
			if err := json.Unmarshal(v, &p[n]); err != nil {
				// TODO: Wrap the error in an InvalidPage?
				return err
			}
			n++
			k, v = cursor.Next()
			if k == nil {
				break
			}
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadStateKey{SeqNum: c, BadKey: key}
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

// CreatePage creates a new page within a collection.
func (repo *Repository) CreatePage(c collection.SeqNum, pageNum page.SeqNum, data page.Data) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return repo.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
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
				return BadPageKey{SeqNum: c, BadKey: key}
			default:
				expected = page.SeqNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if pageNum != expected {
			return collection.PageOutOfOrder{SeqNum: c, PageNum: pageNum, Expected: expected}
		}

		key := makePageKey(pageNum)
		return pages.Put(key[:], value)
	})
}

// ClearPages removes pages affiliated with a collection.
func (repo *Repository) ClearPages(c collection.SeqNum) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		col := collectionBucket(tx, repo.teamDriveID, c)
		if col == nil {
			return collection.NotFound{SeqNum: c}
		}
		key := []byte(PageBucket)
		pages := col.Bucket(key)
		if pages == nil {
			return nil
		}
		return col.DeleteBucket(key)
	})
}
