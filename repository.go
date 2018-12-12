package drivestream

import "github.com/scjalliance/drivestream/resource"

// Repository is an interface that provides access to drivestream data
// for a particular team drive.
type Repository interface {
	// Type returns a string describing the type of the repository.
	Type() string

	// Drives returns a drive map.
	Drives() DriveMap

	// Drive returns a drive reference.
	Drive(driveID resource.ID) DriveReference

	// Files returns a file map.
	Files() FileMap

	// File returns a file reference.
	File(fileID resource.ID) FileReference
}
