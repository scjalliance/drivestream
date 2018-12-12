package boltrepo

import (
	"fmt"

	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// BadCollectionKey reports that the repository contains invalid key
// data within its collection table.
type BadCollectionKey struct {
	Drive  resource.ID
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCollectionKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: the database contains an invalid collection key: %v", e.Drive, e.BadKey)
}

// BadCollectionStateKey reports that a collection contains contains an
// invalid state key.
type BadCollectionStateKey struct {
	Drive      resource.ID
	Collection collection.SeqNum
	BadKey     []byte
}

// Error returns a string representation of the error.
func (e BadCollectionStateKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d: the database contains an invalid state key: %v", e.Drive, e.Collection, e.BadKey)
}

// BadPageKey reports that a collection contains an invalid page key.
type BadPageKey struct {
	Drive      resource.ID
	Collection collection.SeqNum
	BadKey     []byte
}

// Error returns a string representation of the error.
func (e BadPageKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d: the database contains an invalid page key: %v", e.Drive, e.Collection, e.BadKey)
}

// BadCommitKey reports that the repository contains invalid key
// data within its commit table.
type BadCommitKey struct {
	Drive  resource.ID
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCommitKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: the database contains an invalid commit key: %v", e.Drive, e.BadKey)
}

// BadCommitStateKey reports that a commit contains contains an
// invalid state key.
type BadCommitStateKey struct {
	Drive  resource.ID
	Commit commit.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCommitStateKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d: the database contains an invalid state key: %v", e.Drive, e.Commit, e.BadKey)
}

// BadCommitFileVersion reports that a commit contains an invalid file
// version key.
type BadCommitFileVersion struct {
	Drive  resource.ID
	Commit commit.SeqNum
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadCommitFileVersion) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d: the database contains an invalid file version key: %v", e.Drive, e.Commit, e.BadKey)
}

// BadDriveVersionKey reports that the repository contains invalid key
// data within its drive version table.
type BadDriveVersionKey struct {
	Drive  resource.ID
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadDriveVersionKey) Error() string {
	return fmt.Sprintf("drivestream: drive %s: the database contains an invalid drive version key: %v", e.Drive, e.BadKey)
}

// BadFileVersionKey reports that the repository contains invalid key
// data within its file version table.
type BadFileVersionKey struct {
	File   resource.ID
	BadKey []byte
}

// Error returns a string representation of the error.
func (e BadFileVersionKey) Error() string {
	return fmt.Sprintf("drivestream: file %s: the database contains an invalid drive version key: %v", e.File, e.BadKey)
}
