package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.StateSequence = (*CommitStates)(nil)

// CommitStates accesses a sequence of commit states in a
// bolt repository.
type CommitStates struct {
	db     *bolt.DB
	drive  resource.ID
	commit commit.SeqNum
}

// Next returns the state number to use for the next state.
func (seq CommitStates) Next() (n commit.StateNum, err error) {
	err = seq.db.View(func(tx *bolt.Tx) error {
		com := commitBucket(tx, seq.drive, seq.commit)
		if com == nil {
			return commit.NotFound{Drive: seq.drive, Commit: seq.commit}
		}
		states := com.Bucket([]byte(StateBucket))
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
			return BadCommitStateKey{Drive: seq.drive, Commit: seq.commit, BadKey: key}
		}
		n = commit.StateNum(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads a subset of states from the sequence, starting at start.
// Up to len(p) states will be returned in p. The number of states
// returned is provided as n.
func (seq CommitStates) Read(start commit.StateNum, p []commit.State) (n int, err error) {
	err = seq.db.View(func(tx *bolt.Tx) error {
		com := commitBucket(tx, seq.drive, seq.commit)
		if com == nil {
			return commit.NotFound{Drive: seq.drive, Commit: seq.commit}
		}
		states := com.Bucket([]byte(StateBucket))
		if states == nil {
			return commit.StateNotFound{Drive: seq.drive, Commit: seq.commit, State: start}
		}
		cursor := states.Cursor()
		pos := start
		key := makeCommitStateKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return commit.StateNotFound{Drive: seq.drive, Commit: seq.commit, State: start}
		}
		for n < len(p) {
			if v == nil {
				return commit.StateInvalid{Drive: seq.drive, Commit: seq.commit, State: pos} // All states must be non-nil
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
				return BadCommitStateKey{Drive: seq.drive, Commit: seq.commit, BadKey: key}
			}
			pos = start + commit.StateNum(n)
			key = makeCommitStateKey(pos)
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

// Ref returns a commit state reference for the sequence number.
func (seq CommitStates) Ref(stateNum commit.StateNum) commit.StateReference {
	return CommitState{
		db:     seq.db,
		drive:  seq.drive,
		commit: seq.commit,
		state:  stateNum,
	}
}
