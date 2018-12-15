package boltrepo

import (
	"bytes"
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ fileview.Reference = (*FileView)(nil)

// FileView is a drivestream file version reference for a bolt
// repository.
type FileView struct {
	db    *bolt.DB
	file  resource.ID
	drive resource.ID
}

// File returns the ID of the file.
func (ref FileView) File() resource.ID {
	return ref.file
}

// Drive returns the ID of the drive being viewed.
func (ref FileView) Drive() resource.ID {
	return ref.drive
}

// At returns the version reference of the file at a particular commit.
//
// TODO: Consider returning the closest commit number as well as the version.
func (ref FileView) At(seqNum commit.SeqNum) (r fileversion.Reference, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		views := fileViewsBucket(tx, ref.file)
		if views == nil {
			return fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
		}
		view := views.Bucket([]byte(ref.drive))
		if view == nil {
			return fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
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
			// The file didn't exit within the drive at seqNum.
			return fileview.NotFound{File: ref.file, Drive: ref.drive, Commit: seqNum}
		}

		if len(k) != 8 {
			key := append(k[:0:0], k...) // Copy key bytes
			return BadFileViewKey{File: ref.file, Drive: ref.drive, BadKey: key}
		}

		if len(v) != 8 {
			value := append(v[:0:0], v...) // Copy value bytes
			return BadFileViewValue{File: ref.file, Drive: ref.drive, Commit: commit.SeqNum(binary.BigEndian.Uint64(k)), BadValue: value}
		}

		// FIXME: What about deleted files?
		// TODO: Include deletions in the view with a -1 version number so we
		//       don't pick up the pre-deleted states?

		r = FileVersion{
			db:      ref.db,
			file:    ref.file,
			version: resource.Version(binary.BigEndian.Uint64(v)),
		}
		return nil
	})
	return r, err
}

// Add adds version as a view of the file at the commit sequence number.
func (ref FileView) Add(seqNum commit.SeqNum, version resource.Version) error {
	return ref.db.Update(func(tx *bolt.Tx) error {
		views, err := createFileViewsBucket(tx, ref.file)
		if err != nil {
			return err
		}

		view, err := views.CreateBucketIfNotExists([]byte(ref.drive))
		if err != nil {
			return err
		}

		key := makeCommitKey(seqNum)
		value := makeVersionKey(version)
		return view.Put(key[:], value[:])
	})
}
