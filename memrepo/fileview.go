package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileview.Reference = (*FileView)(nil)

// FileView is a drivestream file version reference for a bolt
// repository.
type FileView struct {
	repo  *Repository
	file  resource.ID
	drive resource.ID
}

// File returns the ID of the file.
func (ref FileView) File() resource.ID {
	return ref.file
}

// Drive returns the ID of the drive being viewed.
func (ref FileView) Drive() resource.ID {
	return ref.drive
}

// At returns the version reference of the file at a particular commit.
//
// TODO: Consider returning the closest commit number as well as the version.
func (ref FileView) At(seqNum commit.SeqNum) (r fileversion.Reference, err error) {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		return nil, fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
	}
	view, ok := file.Views[ref.drive]
	if !ok {
		return nil, fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
	}
	var (
		maxCommit  commit.SeqNum
		maxVersion resource.Version
		found      bool
	)
	for seqNum, version := range view {
		if !found || seqNum > maxCommit {
			found = true
			maxCommit = seqNum
			maxVersion = version
		}
	}
	if !found {
		return nil, fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
	}
	// FIXME: What about deleted files?
	// TODO: Include deletions in the view with a -1 version number so we
	//       don't pick up the pre-deleted states?
	return FileVersion{
		repo:    ref.repo,
		file:    ref.file,
		version: maxVersion,
	}, nil
}

// Add adds version as a view of the file at the commit sequence number.
func (ref FileView) Add(seqNum commit.SeqNum, version resource.Version) error {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		file = newFileEntry()
	}
	view, ok := file.Views[ref.drive]
	if !ok {
		view = make(map[commit.SeqNum]resource.Version)
		file.Views[ref.drive] = view
	}
	view[seqNum] = version
	ref.repo.files[ref.file] = file
	return nil
}
