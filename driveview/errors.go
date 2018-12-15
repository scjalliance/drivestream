package driveview

import (
	"fmt"

	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a view of the drive could not be found within
// the repository. This typically means that commit 0 for the drive hasn't
// been finalized.
type NotFound struct {
	Drive  resource.ID
	Commit commit.SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: view not found for commit %d", e.Drive, e.Commit)
}
