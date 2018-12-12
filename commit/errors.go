package commit

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a commit could not be found
// within the repository.
type NotFound struct {
	Drive  resource.ID
	Commit SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d could not be found", e.Drive, e.Commit)
}

// OutOfOrder reports that a commit could not be created
// because its sequence number is not the next one in the series.
type OutOfOrder struct {
	Drive    resource.ID
	Commit   SeqNum
	Expected SeqNum
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d could not be created because it would be out of order (expected %d)", e.Drive, e.Commit, e.Expected)
}

// DataInvalid reports that a commit contains invalid or
// unparsable data.
type DataInvalid struct {
	Drive  resource.ID
	Commit SeqNum
}

// Error returns a string representation of the error.
func (e DataInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d contains invalid data", e.Drive, e.Commit)
}

// StateNotFound reports that a requested state does not exist
// within the commit.
type StateNotFound struct {
	Drive  resource.ID
	Commit SeqNum
	State  StateNum
}

// Error returns a string representation of the error.
func (e StateNotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d does not contain state %d", e.Drive, e.Commit, e.State)
}

// StateOutOfOrder reports that a commit state could not be created
// because its state number is not the next one in the series.
type StateOutOfOrder struct {
	Drive    resource.ID
	Commit   SeqNum
	State    StateNum
	Expected StateNum
}

// Error returns a string representation of the error.
func (e StateOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d state %d could not be created because it would be out of order (expected %d)", e.Drive, e.Commit, e.State, e.Expected)
}

// StateInvalid reports that a requested commit contains invalid
// or unparsable state.
type StateInvalid struct {
	Drive  resource.ID
	Commit SeqNum
	State  StateNum
}

// Error returns a string representation of the error.
func (e StateInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d contains invalid data for state %d", e.Drive, e.Commit, e.State)
}

// StatesTruncated reports that a request for commit states
// returned a shorter list than expected. This is an indication that the
// repository is providing an inconsistent view of its data.
type StatesTruncated struct {
	Drive  resource.ID
	Commit SeqNum
}

// Error returns a string representation of the error.
func (e StatesTruncated) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d contains an inconsistent view of its states", e.Drive, e.Commit)
}

// TreeGroupNotFound reports that a requested group of tree changes does not
// exist within the commit.
type TreeGroupNotFound struct {
	Drive  resource.ID
	Commit SeqNum
	Parent resource.ID
}

// Error returns a string representation of the error.
func (e TreeGroupNotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d does not contain tree change group %s", e.Drive, e.Commit, e.Parent)
}

// TreeGroupInvalid reports that a commit contains invalid data in a group of
// tree changes.
type TreeGroupInvalid struct {
	Drive  resource.ID
	Commit SeqNum
	Parent resource.ID
}

// Error returns a string representation of the error.
func (e TreeGroupInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: commit %d contains invalid tree change data in group %s", e.Drive, e.Commit, e.Parent)
}
