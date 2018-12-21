package memrepo

import (
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileMap = (*Files)(nil)

// Files accesses a map of files in a bolt repository.
type Files struct {
	repo *Repository
}

// List returns the list of files contained within the repository.
func (ref Files) List() (ids []resource.ID, err error) {
	for id := range ref.repo.files {
		ids = append(ids, id)
	}
	return ids, nil
}

// Ref returns a file reference.
func (ref Files) Ref(id resource.ID) drivestream.FileReference {
	return File{
		repo: ref.repo,
		file: id,
	}
}

// AddVersions adds file versions to the file map in bulk.
func (ref Files) AddVersions(fileVersions ...resource.File) error {
	for _, fileVersion := range fileVersions {
		file, ok := ref.repo.files[fileVersion.ID]
		if !ok {
			file = newFileEntry()
		}
		file.Versions[fileVersion.Version] = fileVersion.FileData
		ref.repo.files[fileVersion.ID] = file
	}
	return nil
}

// AddViewData adds view data to the file map in bulk.
func (ref Files) AddViewData(entries ...fileview.Data) error {
	for _, entry := range entries {
		file, ok := ref.repo.files[entry.File]
		if !ok {
			file = newFileEntry()
		}
		view, ok := file.Views[entry.Drive]
		if !ok {
			view = make(map[commit.SeqNum]resource.Version)
			file.Views[entry.Drive] = view
		}
		view[entry.Commit] = entry.Version
		ref.repo.files[entry.File] = file
	}
	return nil
}
