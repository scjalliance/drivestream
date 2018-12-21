package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.Sequence = (*Commits)(nil)

// Commits accesses a sequence of commits in a bolt repository.
type Commits struct {
	db    *bolt.DB
	drive resource.ID
}

// Next returns the sequence number to use for the next commit.
func (ref Commits) Next() (n commit.SeqNum, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		commits := commitsBucket(tx, ref.drive)
		if commits == nil {
			return nil
		}
		cursor := commits.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return nil
		}
		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadCommitKey{Drive: ref.drive, BadKey: key}
		}
		n = commit.SeqNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads commit data for a range of commits
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (ref Commits) Read(start commit.SeqNum, p []commit.Data) (n int, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		commits := commitsBucket(tx, ref.drive)
		if commits == nil {
			return commit.NotFound{Drive: ref.drive, Commit: start}
		}
		cursor := commits.Cursor()
		pos := start
		key := makeCommitKey(pos)
		k, _ := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return commit.NotFound{Drive: ref.drive, Commit: start}
		}
		for n < len(p) {
			col := commits.Bucket(k)
			if col == nil {
				return commit.DataInvalid{Drive: ref.drive, Commit: pos} // All commits must be buckets
			}
			value := col.Get([]byte(DataKey))
			if value == nil {
				return commit.DataInvalid{Drive: ref.drive, Commit: pos}
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
			pos = start + commit.SeqNum(n)
			key = makeCommitKey(pos)
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

// Ref returns a commit reference.
func (ref Commits) Ref(c commit.SeqNum) commit.Reference {
	return Commit{
		db:     ref.db,
		drive:  ref.drive,
		commit: c,
	}
}
