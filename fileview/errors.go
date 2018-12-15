package fileview

import (
	"fmt"

	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a view of a file could not be found within the
// repository for the requested drive and commit.
type NotFound struct {
	File   resource.ID
	Drive  resource.ID
	Commit commit.SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: file %s: view not found for drive %s commit %d", e.File, e.Drive, e.Commit)
}
