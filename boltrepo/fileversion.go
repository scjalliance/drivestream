package boltrepo

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileversion.Reference = (*FileVersion)(nil)

// FileVersion is a drivestream file version reference for a bolt
// repository.
type FileVersion struct {
	db      *bolt.DB
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
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		versions, err := createFileVersionsBucket(tx, ref.file)
		if err != nil {
			return err
		}

		key := makeVersionKey(ref.version)
		return versions.Put(key[:], value)
	})
}

// Data returns the data of the file version.
func (ref FileVersion) Data() (data resource.FileData, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		versions := fileVersionsBucket(tx, ref.file)
		if versions == nil {
			return fileversion.NotFound{File: ref.file, Version: ref.version}
		}
		key := makeVersionKey(ref.version)
		value := versions.Get(key[:])
		if value == nil {
			return fileversion.NotFound{File: ref.file, Version: ref.version}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}
