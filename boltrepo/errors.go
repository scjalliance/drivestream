package boltrepo

import (
	"fmt"

	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
)

// BadCollectionKey reports that the repository contains invalid key
// data within its collection table.
type BadCollectionKey struct {
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCollectionKey) Error() string {
	return fmt.Sprintf("the database contains an invalid collection key: %v", e.BadKey)
}

// BadCollectionStateKey reports that a collection contains contains an
// invalid state key.
type BadCollectionStateKey struct {
	SeqNum collection.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCollectionStateKey) Error() string {
	return fmt.Sprintf("the database contains an invalid state key within collection %d: %v", e.SeqNum, e.BadKey)
}

// BadPageKey reports that a collection contains an invalid page key.
type BadPageKey struct {
	SeqNum collection.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadPageKey) Error() string {
	return fmt.Sprintf("the database contains an invalid page key within collection %d: %v", e.SeqNum, e.BadKey)
}

// BadCommitKey reports that the repository contains invalid key
// data within its commit table.
type BadCommitKey struct {
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCommitKey) Error() string {
	return fmt.Sprintf("the database contains an invalid commit key: %v", e.BadKey)
}

// BadCommitStateKey reports that a commit contains contains an
// invalid state key.
type BadCommitStateKey struct {
	SeqNum commit.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCommitStateKey) Error() string {
	return fmt.Sprintf("the database contains an invalid state key within commit %d: %v", e.SeqNum, e.BadKey)
}
