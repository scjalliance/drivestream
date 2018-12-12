package boltrepo

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

// FileVersions accesses a map of file versions in a bolt repository.
type FileVersions struct {
	db   *bolt.DB
	file resource.ID
}

// List returns a list of version numbers for the file.
func (versionMap FileVersions) List() (v []resource.Version, err error) {
	err = versionMap.db.View(func(tx *bolt.Tx) error {
		versions := fileVersionsBucket(tx, versionMap.file)
		if versions == nil {
			return nil
		}
		cursor := versions.Cursor()
		k, _ := cursor.First()
		for k != nil {
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadFileVersionKey{File: versionMap.file, BadKey: key}
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
func (versionMap FileVersions) Ref(v resource.Version) fileversion.Reference {
	return FileVersion{
		db:      versionMap.db,
		file:    versionMap.file,
		version: v,
	}
}
