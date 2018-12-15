package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/driveview"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveview.Reference = (*DriveView)(nil)

// DriveView is a drivestream drive version reference for a bolt
// repository.
type DriveView struct {
	repo  *Repository
	drive resource.ID
}

// Drive returns the ID of the drive being viewed.
func (ref DriveView) Drive() resource.ID {
	return ref.drive
}

// At returns the version reference of the drive at a particular commit.
//
// TODO: Consider returning the closest commit number as well as the version.
func (ref DriveView) At(seqNum commit.SeqNum) (r driveversion.Reference, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return nil, driveview.NotFound{Drive: ref.drive, Commit: seqNum}
	}
	var (
		maxCommit  commit.SeqNum
		maxVersion resource.Version
		found      bool
	)
	for seqNum, version := range drv.View {
		if !found || seqNum > maxCommit {
			found = true
			maxCommit = seqNum
			maxVersion = version
		}
	}
	if !found {
		return nil, driveview.NotFound{Drive: ref.drive, Commit: seqNum}
	}
	return DriveVersion{
		repo:    ref.repo,
		drive:   ref.drive,
		version: maxVersion,
	}, nil
}

// Add adds version as a view of the drive at the commit sequence number.
func (ref DriveView) Add(seqNum commit.SeqNum, version resource.Version) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		drv = newDriveEntry()
	}
	drv.View[seqNum] = version
	ref.repo.drives[ref.drive] = drv
	return nil
}
