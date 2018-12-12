package driveversion

import (
	"github.com/scjalliance/drivestream/resource"
)

// Reference is a drive version reference.
type Reference interface {
	// Drive returns the ID of the drive.
	Drive() resource.ID

	// Version returns the version number of the drive.
	Version() resource.Version

	// Create creates a new drive version with the given version number
	// and data. If a version already exists with the version number an
	// error will be returned.
	Create(data resource.DriveData) error

	// Data returns the data of the drive version.
	Data() (resource.DriveData, error)
}
