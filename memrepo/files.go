package memrepo

import (
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileMap = (*Files)(nil)

// Files accesses a map of files in a bolt repository.
type Files struct {
	repo *Repository
}

// List returns the list of files contained within the repository.
func (fileMap Files) List() (ids []resource.ID, err error) {
	for id := range fileMap.repo.files {
		ids = append(ids, id)
	}
	return ids, nil
}

// Ref returns a file reference.
func (fileMap Files) Ref(id resource.ID) drivestream.FileReference {
	return File{
		repo: fileMap.repo,
		file: id,
	}
}

// AddVersions adds file versions to the file map in bulk.
func (fileMap Files) AddVersions(fileVersions ...resource.File) error {
	for _, fileVersion := range fileVersions {
		file, ok := fileMap.repo.files[fileVersion.ID]
		if !ok {
			file = newFileEntry()
		}
		file.Versions[fileVersion.Version] = fileVersion.FileData
		fileMap.repo.files[fileVersion.ID] = file
	}
	return nil
}
