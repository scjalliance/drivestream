package drivestream

import "github.com/scjalliance/drivestream/resource"

// DriveMap is a map of drivestream drives.
type DriveMap interface {
	// List returns a list of all drives within the map.
	List() (ids []resource.ID, err error)

	// Ref returns a drive reference.
	Ref(driveID resource.ID) DriveReference
}
