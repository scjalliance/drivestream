package boltrepo

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileversion.Map = (*FileVersions)(nil)

// FileVersions accesses a map of file versions in a bolt repository.
type FileVersions struct {
	db   *bolt.DB
	file resource.ID
}

// Path returns the path of the file versions.
func (ref FileVersions) Path() binpath.Text {
	return binpath.Text{RootBucket, FileBucket, ref.file.String(), VersionBucket}
}

// List returns a list of version numbers for the file.
func (ref FileVersions) List() (v []resource.Version, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		versions := fileVersionsBucket(tx, ref.file)
		if versions == nil {
			return nil
		}
		cursor := versions.Cursor()
		k, _ := cursor.First()
		for k != nil {
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadFileVersionKey{File: ref.file, BadKey: key}
			}
			version := resource.Version(binary.BigEndian.Uint64(k))
			v = append(v, version)
			k, _ = cursor.Next()
		}
		return nil
	})
	return v, err
}

// Ref returns a file version reference for the version number.
func (ref FileVersions) Ref(v resource.Version) fileversion.Reference {
	return FileVersion{
		db:      ref.db,
		file:    ref.file,
		version: v,
	}
}
