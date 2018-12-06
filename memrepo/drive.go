package memrepo

import (
	"github.com/scjalliance/drivestream/resource"
)

// Drive holds version and commit history for a team drive.
type Drive struct {
	ID       resource.ID
	Versions []resource.DriveData
	//Commits  map[commit.SeqNum]drivestream.Version
}
