package drivestream

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/driveview"
	"github.com/scjalliance/drivestream/resource"
)

// DriveReference is a reference to a drivestream drive.
type DriveReference interface {
	// DriveID returns the resource ID of the drive.
	DriveID() resource.ID

	// Exists returns true if the drive exists.
	Exists() (bool, error)

	// Collections returns the collection sequence for the drive.
	Collections() collection.Sequence

	// Collection returns a collection reference. Equivalent to Collections().Ref(s).
	Collection(s collection.SeqNum) collection.Reference

	// Commits returns the commit sequence for the drive.
	Commits() commit.Sequence

	// Commit returns a commit reference. Equivalent to Commits().Ref(s).
	Commit(s commit.SeqNum) commit.Reference

	// Versions returns the version sequence for the drive.
	Versions() driveversion.Sequence

	// Version returns a drive version reference. Equivalent to Versions().Ref(s).
	Version(v resource.Version) driveversion.Reference

	// View returns a view of the drive.
	View() driveview.Reference

	// At returns a version reference of the drive at a particular commit.
	At(seqNum commit.SeqNum) (driveversion.Reference, error)

	// Tree returns the tree map for the drive.
	//Tree() drivetree.Map

	// Stats returns statistics about the drive.
	Stats() (DriveStats, error)
}
