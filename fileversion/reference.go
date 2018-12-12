package fileversion

import "github.com/scjalliance/drivestream/resource"

// Reference is a file version reference.
type Reference interface {
	// File returns the ID of the file.
	File() resource.ID

	// Version returns the version number of the file.
	Version() resource.Version

	// Create creates a new file version with the given version number
	// and data. If a version already exists with the version number an
	// error will be returned.
	Create(data resource.FileData) error

	// Data returns the data of the file version.
	Data() (resource.FileData, error)
}
