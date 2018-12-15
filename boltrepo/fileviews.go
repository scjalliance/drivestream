package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileview.Map = (*FileViews)(nil)

// FileViews accesses a map of file views in a bolt repository.
type FileViews struct {
	db   *bolt.DB
	file resource.ID
}

// List returns a list of drives with a view of the file.
func (ref FileViews) List() (drives []resource.ID, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		views := fileViewsBucket(tx, ref.file)
		if views == nil {
			return nil
		}
		cursor := views.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			drives = append(drives, resource.ID(k))
		}
		return nil
	})
	return drives, err
}

// Ref returns a view of the file for a particular drive.
func (ref FileViews) Ref(driveID resource.ID) fileview.Reference {
	return FileView{
		db:    ref.db,
		file:  ref.file,
		drive: driveID,
	}
}
