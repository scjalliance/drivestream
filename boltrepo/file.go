package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileReference = (*File)(nil)

// File is a drivestream file reference for a bolt repository.
type File struct {
	db   *bolt.DB
	file resource.ID
}

// Path returns the path of the file.
func (ref File) Path() binpath.Text {
	return binpath.Text{RootBucket, FileBucket, ref.file.String()}
}

// FileID returns the resource ID of the file.
func (ref File) FileID() resource.ID {
	return ref.file
}

// Exists returns true if the file exists.
func (ref File) Exists() (exists bool, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		if fileBucket(tx, ref.file) != nil {
			exists = true
		}
		return nil
	})
	return exists, err
}

// Versions returns the version map for the file.
func (ref File) Versions() fileversion.Map {
	return FileVersions{
		db:   ref.db,
		file: ref.file,
	}
}

// Version returns a file version reference. Equivalent to Versions().Ref(s).
func (ref File) Version(v resource.Version) fileversion.Reference {
	return FileVersion{
		db:      ref.db,
		file:    ref.file,
		version: v,
	}
}

// Views returns the view map for the file.
func (ref File) Views() fileview.Map {
	return FileViews{
		db:   ref.db,
		file: ref.file,
	}
}

// View returns a view of the file for a particular drive.
func (ref File) View(driveID resource.ID) fileview.Reference {
	return FileView{
		db:    ref.db,
		file:  ref.file,
		drive: driveID,
	}
}
