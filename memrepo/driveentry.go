package memrepo

import "github.com/scjalliance/drivestream/resource"

// DriveEntry holds version and commit history for a drive.
type DriveEntry struct {
	Collections []CollectionEntry
	Commits     []CommitEntry
	Versions    []resource.DriveData
}

func newDriveEntry() DriveEntry {
	return DriveEntry{}
}
