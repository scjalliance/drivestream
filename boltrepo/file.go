package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileReference = (*File)(nil)

// File is a drivestream file reference for a bolt repository.
type File struct {
	db   *bolt.DB
	file resource.ID
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
