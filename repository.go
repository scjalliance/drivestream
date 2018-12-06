package drivestream

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

// Repository is an interface that provides access to drivestream data
// for a particular team drive.
type Repository interface {
	// DriveID returns the team drive ID.
	DriveID() resource.ID

	collection.Repository
	//commit.Repository
	//filetree.Repository
}
