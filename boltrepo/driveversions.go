package boltrepo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveversion.Sequence = (*DriveVersions)(nil)

// DriveVersions accesses a sequence of drive versions in a bolt repository.
type DriveVersions struct {
	db    *bolt.DB
	drive resource.ID
}

// Next returns the next version number in the sequence.
func (ref DriveVersions) Next() (n resource.Version, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		versions := driveVersionsBucket(tx, ref.drive)
		if versions == nil {
			return nil
		}
		cursor := versions.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return nil
		}
		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadDriveVersionKey{Drive: ref.drive, BadKey: key}
		}
		n = resource.Version(binary.BigEndian.Uint64(k)) + 1
		return nil
	})
	return n, err
}

// Read reads drive data for a range of drive versions starting at the
// given version number. Up to len(p) entries will be returned in p.
// The number of entries is returned as n.
func (ref DriveVersions) Read(start resource.Version, p []resource.DriveData) (n int, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		versions := driveVersionsBucket(tx, ref.drive)
		if versions == nil {
			return driveversion.NotFound{Drive: ref.drive, Version: start}
		}
		cursor := versions.Cursor()
		pos := start
		key := makeVersionKey(pos)
		k, v := cursor.Seek(key[:])
		if k == nil || !bytes.Equal(key[:], k) {
			return driveversion.NotFound{Drive: ref.drive, Version: start}
		}
		for n < len(p) {
			if v == nil {
				return driveversion.InvalidData{Drive: ref.drive, Version: pos} // All versions must be non-nil
			}
			if err := json.Unmarshal(v, &p[n]); err != nil {
				// TODO: Wrap the error in InvalidData?
				return err
			}
			n++
			k, v = cursor.Next()
			if k == nil {
				break
			}
			if len(k) != 8 {
				key := append(k[:0:0], k...) // Copy key bytes
				return BadDriveVersionKey{Drive: ref.drive, BadKey: key}
			}
			pos = start + resource.Version(n)
			key = makeVersionKey(pos)
			if !bytes.Equal(key[:], k) {
				// The next key doesn't match the expected sequence number
				// TODO: Consider returning an error here?
				break
			}
		}
		return nil
	})
	return n, err
}

// Ref returns a drive version reference for the version number.
func (ref DriveVersions) Ref(v resource.Version) driveversion.Reference {
	return DriveVersion{
		db:      ref.db,
		drive:   ref.drive,
		version: v,
	}
}
