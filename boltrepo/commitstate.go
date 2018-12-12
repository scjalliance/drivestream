package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.StateReference = (*CommitState)(nil)

// CommitState is a reference to a commit state.
type CommitState struct {
	db     *bolt.DB
	drive  resource.ID
	commit commit.SeqNum
	state  commit.StateNum
}

// StateNum returns the state number of the reference.
func (ref CommitState) StateNum() commit.StateNum {
	return ref.state
}

// Create creates the commit state with the given data. If a state already
// exists with the state number an error will be returned.
func (ref CommitState) Create(data commit.State) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		states, err := com.CreateBucketIfNotExists([]byte(StateBucket))
		if err != nil {
			return err
		}

		var expected commit.StateNum
		{
			cursor := states.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadCommitStateKey{Drive: ref.drive, Commit: ref.commit, BadKey: key}
			default:
				expected = commit.StateNum(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.state != expected {
			return commit.StateOutOfOrder{Drive: ref.drive, Commit: ref.commit, State: ref.state, Expected: expected}
		}

		key := makeCommitStateKey(ref.state)
		return states.Put(key[:], value)
	})
}

// Data returns the commit state data.
func (ref CommitState) Data() (data commit.State, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		col := commitBucket(tx, ref.drive, ref.commit)
		if col == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		states := col.Bucket([]byte(StateBucket))
		if states == nil {
			return commit.StateNotFound{Drive: ref.drive, Commit: ref.commit, State: ref.state}
		}
		key := makeCommitStateKey(ref.state)
		value := states.Get(key[:])
		if value == nil {
			return commit.StateNotFound{Drive: ref.drive, Commit: ref.commit, State: ref.state}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}
