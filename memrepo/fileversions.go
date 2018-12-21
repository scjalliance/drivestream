package memrepo

import (
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileversion.Map = (*FileVersions)(nil)

// FileVersions accesses a map of file versions in a bolt repository.
type FileVersions struct {
	repo *Repository
	file resource.ID
}

// List returns a list of version numbers for the file.
func (ref FileVersions) List() (v []resource.Version, err error) {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		return nil, nil
	}
	for version := range file.Versions {
		v = append(v, version)
	}
	return v, nil
}

// Ref returns a file version reference for the version number.
func (ref FileVersions) Ref(v resource.Version) fileversion.Reference {
	return FileVersion{
		repo:    ref.repo,
		file:    ref.file,
		version: v,
	}
}
