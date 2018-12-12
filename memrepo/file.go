package memrepo

import (
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileReference = (*File)(nil)

// File is a drivestream file reference for an in-memory repository.
type File struct {
	repo *Repository
	file resource.ID
}

// FileID returns the resource ID of the file.
func (ref File) FileID() resource.ID {
	return ref.file
}

// Exists returns true if the file exists.
func (ref File) Exists() (exists bool, err error) {
	_, exists = ref.repo.files[ref.file]
	return exists, nil
}

// Versions returns the version map for the file.
func (ref File) Versions() fileversion.Map {
	return FileVersions{
		repo: ref.repo,
		file: ref.file,
	}
}

// Version returns a file version reference. Equivalent to Versions().Ref(s).
func (ref File) Version(v resource.Version) fileversion.Reference {
	return FileVersion{
		repo:    ref.repo,
		file:    ref.file,
		version: v,
	}
}
