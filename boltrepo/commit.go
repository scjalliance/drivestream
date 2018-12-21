package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.Reference = (*Commit)(nil)

// Commit is a drivestream commit reference for a bolt repository.
type Commit struct {
	db     *bolt.DB
	drive  resource.ID
	commit commit.SeqNum
}

// Path returns the path of the commit.
func (ref Commit) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, ref.drive.String(), CommitBucket, ref.commit.String()}
}

// Drive returns the drive ID of the commit.
func (ref Commit) Drive() resource.ID {
	return ref.drive
}

// SeqNum returns the sequence number of the commit.
func (ref Commit) SeqNum() commit.SeqNum {
	return ref.commit
}

// Exists returns true if the commit exists.
func (ref Commit) Exists() (exists bool, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		if commitBucket(tx, ref.drive, ref.commit) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// Create creates a new commit with the given sequence number and data.
// If a commit already exists with the sequence number an error will be
// returned.
func (ref Commit) Create(data commit.Data) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		commits, err := createCommitsBucket(tx, ref.drive)
		if err != nil {
			return err
		}

		var expected commit.SeqNum
		{
			cursor := commits.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCommitKey{Drive: ref.drive, BadKey: key}
			default:
				expected = commit.SeqNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.commit != expected {
			return commit.OutOfOrder{Drive: ref.drive, Commit: ref.commit, Expected: expected}
		}

		key := makeCommitKey(ref.commit)
		col, err := commits.CreateBucket(key[:])
		if err != nil {
			return err
		}
		return col.Put([]byte(DataKey), value)
	})
}

// Data returns information about the commit.
func (ref Commit) Data() (data commit.Data, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := commitBucket(tx, ref.drive, ref.commit)
		if col == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		value := col.Get([]byte(DataKey))
		if value == nil {
			return commit.DataInvalid{Drive: ref.drive, Commit: ref.commit}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}

// States returns the state sequence for the commit.
func (ref Commit) States() commit.StateSequence {
	return CommitStates{
		db:     ref.db,
		drive:  ref.drive,
		commit: ref.commit,
	}
}

// State returns a state reference.
func (ref Commit) State(stateNum commit.StateNum) commit.StateReference {
	return ref.States().Ref(stateNum)
}

// Files returns the map of file changes for the commit.
func (ref Commit) Files() commit.FileMap {
	return CommitFiles{
		db:     ref.db,
		drive:  ref.drive,
		commit: ref.commit,
	}
}

// Tree returns the map of tree changes for the commit.
func (ref Commit) Tree() commit.TreeMap {
	return CommitTree{
		db:     ref.db,
		drive:  ref.drive,
		commit: ref.commit,
	}
}
