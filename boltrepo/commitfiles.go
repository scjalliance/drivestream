package boltrepo

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.FileMap = (*CommitFiles)(nil)

// CommitFiles is a reference to a commit file map.
type CommitFiles struct {
	db     *bolt.DB
	drive  resource.ID
	commit commit.SeqNum
}

// Read returns the set of file changes for the commit, in unspecified
// order.
func (ref CommitFiles) Read() (changes []commit.FileChange, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		files := com.Bucket([]byte(FileBucket))
		if files == nil {
			return nil
		}
		cursor := files.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if len(v) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCommitFileVersion{Drive: ref.drive, Commit: ref.commit, BadKey: key}
			}
			changes = append(changes, commit.FileChange{
				File:    resource.ID(k),
				Version: resource.Version(binary.BigEndian.Uint64(v)),
			})
		}
		return nil
	})
	return changes, err
}

// Add adds the given file changes to the map.
// If two or more changes conflict, the last change added takes
// precedence.
func (ref CommitFiles) Add(changes ...commit.FileChange) error {
	return ref.db.Update(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		files, err := com.CreateBucketIfNotExists([]byte(FileBucket))
		if err != nil {
			return err
		}
		for i := range changes {
			key := []byte(changes[i].File)
			value := makeVersionKey(changes[i].Version)
			if err := files.Put(key, value[:]); err != nil {
				return err
			}
		}
		return nil
	})
}
