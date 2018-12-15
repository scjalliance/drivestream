package memrepo

import (
	"sort"

	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileview.Map = (*FileViews)(nil)

// FileViews accesses a map of file views in a bolt repository.
type FileViews struct {
	repo *Repository
	file resource.ID
}

// List returns a list of drives with a view of the file.
func (ref FileViews) List() (drives []resource.ID, err error) {
	file, ok := ref.repo.files[ref.file]
	if !ok {
		return nil, nil
	}
	drives = make([]resource.ID, 0, len(file.Views))
	for driveID := range file.Views {
		drives = append(drives, driveID)
	}
	sort.Slice(drives, func(i, j int) bool { return drives[i] < drives[j] })
	return drives, err
}

// Ref returns a view of the file for a particular drive.
func (ref FileViews) Ref(driveID resource.ID) fileview.Reference {
	return FileView{
		repo:  ref.repo,
		file:  ref.file,
		drive: driveID,
	}
}
