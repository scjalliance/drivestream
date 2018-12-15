package fileview

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

// Reference is a file view reference.
type Reference interface {
	// File returns the ID of the file.
	File() resource.ID

	// Drive returns the ID of the drive being viewed.
	Drive() resource.ID

	// At returns the version reference of the file at a particular commit.
	At(seqNum commit.SeqNum) (fileversion.Reference, error)

	// Add adds version as a view of the file at the commit sequence number.
	Add(seqNum commit.SeqNum, version resource.Version) error
}
