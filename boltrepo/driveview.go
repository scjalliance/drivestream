package boltrepo

import (
	"bytes"
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/driveview"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveview.Reference = (*DriveView)(nil)

// DriveView is a drivestream drive version reference for a bolt
// repository.
type DriveView struct {
	db    *bolt.DB
	drive resource.ID
}

// Drive returns the ID of the drive being viewed.
func (ref DriveView) Drive() resource.ID {
	return ref.drive
}

// At returns the version reference of the drive at a particular commit.
//
// TODO: Consider returning the closest commit number as well as the version.
func (ref DriveView) At(seqNum commit.SeqNum) (r driveversion.Reference, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		view := driveViewBucket(tx, ref.drive)
		if view == nil {
			return driveview.NotFound{Drive: ref.drive, Commit: seqNum}
		}
		cursor := view.Cursor()
		key := makeCommitKey(seqNum)
		k, v := cursor.Seek(key[:])
		if k == nil {
			// The cursor found no view at or after seqNum.
			k, v = cursor.Last() // Use whatever commit is last
		} else if !bytes.Equal(k, key[:]) {
			// The cursor didn't find seqNum, but found something after it.
			k, v = cursor.Prev() // Back up one commit to whatever came before seqNum
		}

		if k == nil {
			// The drive didn't exist at seqNum.
			// Theoretically this shouldn't be possible unless the initial
			// collection hasn't finished.
			return driveview.NotFound{Drive: ref.drive, Commit: seqNum}
		}

		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadDriveViewKey{Drive: ref.drive, BadKey: key}
		}

		if len(v) != 8 {
			value := append(v[:0:0], v...) // Copy value bytes
			return BadDriveViewValue{Drive: ref.drive, Commit: commit.SeqNum(binary.BigEndian.Uint64(k)), BadValue: value}
		}

		r = DriveVersion{
			db:      ref.db,
			drive:   ref.drive,
			version: resource.Version(binary.BigEndian.Uint64(v)),
		}
		return nil
	})
	return r, err
}

// Add adds version as a view of the drive at the commit sequence number.
func (ref DriveView) Add(seqNum commit.SeqNum, version resource.Version) error {
	return ref.db.Update(func(tx *bolt.Tx) error {
		view, err := createDriveViewBucket(tx, ref.drive)
		if err != nil {
			return err
		}

		key := makeCommitKey(seqNum)
		value := makeVersionKey(version)
		return view.Put(key[:], value[:])
	})
}
