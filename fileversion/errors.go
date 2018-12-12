package fileversion

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a file version could not be found within the
// repository.
type NotFound struct {
	File    resource.ID
	Version resource.Version
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: file %s: version %d could not be found", e.File, e.Version)
}

// OutOfOrder reports that a file version could not be created
// because its version number is not the next one in the series.
type OutOfOrder struct {
	File     resource.ID
	Version  resource.Version
	Expected resource.Version
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: file %s: version %d could not be created because it would be out of order (expected %d)", e.File, e.Version, e.Expected)
}

// InvalidData reports that a requested file version contains invalid
// or unparsable data.
type InvalidData struct {
	File    resource.ID
	Version resource.Version
}

// Error returns a string representation of the error.
func (e InvalidData) Error() string {
	return fmt.Sprintf("drivestream: file %s: version %d contains invalid data", e.File, e.Version)
}
