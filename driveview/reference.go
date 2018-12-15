package driveview

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/resource"
)

// Reference is a drive view reference.
type Reference interface {
	// Drive returns the ID of the drive being viewed.
	Drive() resource.ID

	// At returns the version reference of the file at a particular commit.
	At(seqNum commit.SeqNum) (driveversion.Reference, error)

	// Add adds version as a view of the file at the commit sequence number.
	Add(seqNum commit.SeqNum, version resource.Version) error
}
