package memrepo

import (
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

// FileVersions accesses a map of file versions in a bolt repository.
type FileVersions struct {
	repo *Repository
	file resource.ID
}

// List returns a list of version numbers for the file.
func (versionMap FileVersions) List() (v []resource.Version, err error) {
	file, ok := versionMap.repo.files[versionMap.file]
	if !ok {
		return nil, nil
	}
	for version := range file.Versions {
		v = append(v, version)
	}
	return v, nil
}

// Ref returns a file version reference for the version number.
func (versionMap FileVersions) Ref(v resource.Version) fileversion.Reference {
	return FileVersion{
		repo:    versionMap.repo,
		file:    versionMap.file,
		version: v,
	}
}
