package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// DriveEntry holds version and commit history for a drive.
type DriveEntry struct {
	Collections []CollectionEntry
	Commits     []CommitEntry
	Versions    []resource.DriveData
	View        map[commit.SeqNum]resource.Version
}

func newDriveEntry() DriveEntry {
	return DriveEntry{
		View: make(map[commit.SeqNum]resource.Version),
	}
}
