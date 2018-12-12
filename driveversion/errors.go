package driveversion

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a drive version could not be found within the
// repository.
type NotFound struct {
	Drive   resource.ID
	Version resource.Version
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: version %d could not be found", e.Drive, e.Version)
}

// OutOfOrder reports that a drive version could not be created
// because its version number is not the next one in the series.
type OutOfOrder struct {
	Drive    resource.ID
	Version  resource.Version
	Expected resource.Version
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: version %d could not be created because it would be out of order (expected %d)", e.Drive, e.Version, e.Expected)
}

// InvalidData reports that a requested drive version contains invalid
// or unparsable data.
type InvalidData struct {
	Drive   resource.ID
	Version resource.Version
}

// Error returns a string representation of the error.
func (e InvalidData) Error() string {
	return fmt.Sprintf("drivestream: drive %s: version %d contains invalid data", e.Drive, e.Version)
}
