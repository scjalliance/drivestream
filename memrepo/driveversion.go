package memrepo

import (
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveversion.Reference = (*DriveVersion)(nil)

// DriveVersion is a drivestream drive version reference for an in-memory
// repository.
type DriveVersion struct {
	repo    *Repository
	drive   resource.ID
	version resource.Version
}

// Drive returns the ID of the drive.
func (ref DriveVersion) Drive() resource.ID {
	return ref.drive
}

// Version returns the version number of the drive.
func (ref DriveVersion) Version() resource.Version {
	return ref.version
}

// Create creates a new drive version with the given version number
// and data. If a version already exists with the version number an
// error will be returned.
func (ref DriveVersion) Create(data resource.DriveData) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		drv = newDriveEntry()
	}
	expected := resource.Version(len(drv.Versions))
	if ref.version != expected {
		return driveversion.OutOfOrder{Drive: ref.drive, Version: ref.version, Expected: expected}
	}
	drv.Versions = append(drv.Versions, data)
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns the data of the drive version.
func (ref DriveVersion) Data() (data resource.DriveData, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return resource.DriveData{}, driveversion.NotFound{Drive: ref.drive, Version: ref.version}
	}
	if ref.version >= resource.Version(len(drv.Versions)) {
		return resource.DriveData{}, driveversion.NotFound{Drive: ref.drive, Version: ref.version}
	}
	return drv.Versions[ref.version], nil
}
