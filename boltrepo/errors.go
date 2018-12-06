package boltrepo

import (
	"fmt"

	"github.com/scjalliance/drivestream/collection"
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

// BadStateKey reports that a collection contains contains an invalid state key.
type BadStateKey struct {
	SeqNum collection.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadStateKey) Error() string {
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
