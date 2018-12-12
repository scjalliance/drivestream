package memrepo

import (
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileversion.Reference = (*FileVersion)(nil)

// FileVersion is a drivestream file version reference for an in-memory
// repository.
type FileVersion struct {
	repo    *Repository
	file    resource.ID
	version resource.Version
}

// File returns the ID of the file.
func (ref FileVersion) File() resource.ID {
	return ref.file
}

// Version returns the version number of the file.
func (ref FileVersion) Version() resource.Version {
	return ref.version
}

// Create creates a new file version with the given version number
// and data. If a version already exists with the version number an
// error will be returned.
func (ref FileVersion) Create(data resource.FileData) error {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		file = newFileEntry()
	}
	file.Versions[ref.version] = data
	ref.repo.files[ref.file] = file
	return nil
}

// Data returns the data of the file version.
func (ref FileVersion) Data() (data resource.FileData, err error) {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		return resource.FileData{}, fileversion.NotFound{File: ref.file, Version: ref.version}
	}
	fileData, ok := file.Versions[ref.version]
	if !ok {
		return resource.FileData{}, fileversion.NotFound{File: ref.file, Version: ref.version}
	}
	return fileData, nil
}
