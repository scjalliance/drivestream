package boltrepo

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveversion.Reference = (*DriveVersion)(nil)

// DriveVersion is a drivestream drive version reference for a bolt
// repository.
type DriveVersion struct {
	db      *bolt.DB
	drive   resource.ID
	version resource.Version
}

// Drive returns the ID of the drive.
func (ref DriveVersion) Drive() resource.ID {
	return ref.drive
}

// Version returns the version number of the drive.
func (ref DriveVersion) Version() resource.Version {
	return ref.version
}

// Create creates a new drive version with the given version number
// and data. If a version already exists with the version number an
// error will be returned.
func (ref DriveVersion) Create(data resource.DriveData) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ref.db.Update(func(tx *bolt.Tx) error {
		versions, err := createDriveVersionsBucket(tx, ref.drive)
		if err != nil {
			return err
		}

		var expected resource.Version
		{
			cursor := versions.Cursor()
			k, _ := cursor.Last()
			switch {
			case k == nil:
				expected = 0
			case len(k) != 8:
				key := append(k[:0:0], k...) // Copy key bytes
				return BadDriveVersionKey{BadKey: key}
			default:
				expected = resource.Version(binary.BigEndian.Uint64(k)) + 1
			}
		}
		if ref.version != expected {
			return driveversion.OutOfOrder{Drive: ref.drive, Version: ref.version, Expected: expected}
		}

		key := makeVersionKey(ref.version)
		return versions.Put(key[:], value)
	})
}

// Data returns the data of the drive version.
func (ref DriveVersion) Data() (data resource.DriveData, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		versions := driveVersionsBucket(tx, ref.drive)
		if versions == nil {
			return driveversion.NotFound{Drive: ref.drive, Version: ref.version}
		}
		key := makeVersionKey(ref.version)
		value := versions.Get(key[:])
		if value == nil {
			return driveversion.NotFound{Drive: ref.drive, Version: ref.version}
		}
		if err := json.Unmarshal(value, &data); err != nil {
			// TODO: Wrap the error in DataInvalid?
			return err
		}
		return nil
	})
	return data, err
}
