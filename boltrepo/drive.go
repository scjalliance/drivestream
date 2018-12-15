package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/driveview"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.DriveReference = (*Drive)(nil)

// Drive is a drivestream drive reference for a bolt repository.
type Drive struct {
	db    *bolt.DB
	drive resource.ID
}

// DriveID returns the resource ID of the drive.
func (ref Drive) DriveID() resource.ID {
	return ref.drive
}

// Exists returns true if the drive exists.
func (ref Drive) Exists() (exists bool, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		if driveBucket(tx, ref.drive) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// Collections returns the collection sequence for the drive.
func (ref Drive) Collections() collection.Sequence {
	return Collections{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Collection returns a collection reference. Equivalent to Collections().Ref(s).
func (ref Drive) Collection(c collection.SeqNum) collection.Reference {
	return Collection{
		db:         ref.db,
		drive:      ref.drive,
		collection: c,
	}
}

// Commits returns the commit sequence for the drive.
func (ref Drive) Commits() commit.Sequence {
	return Commits{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Commit returns a commit reference. Equivalent to Commits().Ref(s).
func (ref Drive) Commit(c commit.SeqNum) commit.Reference {
	return Commit{
		db:     ref.db,
		drive:  ref.drive,
		commit: c,
	}
}

// Versions returns the version sequence for the drive.
func (ref Drive) Versions() driveversion.Sequence {
	return DriveVersions{
		db:    ref.db,
		drive: ref.drive,
	}
}

// Version returns a drive version reference. Equivalent to Versions().Ref(s).
func (ref Drive) Version(v resource.Version) driveversion.Reference {
	return DriveVersion{
		db:      ref.db,
		drive:   ref.drive,
		version: v,
	}
}

// View returns a view of the drive.
func (ref Drive) View() driveview.Reference {
	return DriveView{
		db:    ref.db,
		drive: ref.drive,
	}
}

// At returns a version reference of the drive at a particular commit.
func (ref Drive) At(seqNum commit.SeqNum) (driveversion.Reference, error) {
	return ref.View().At(seqNum)
}
