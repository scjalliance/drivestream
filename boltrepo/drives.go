package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.DriveMap = (*Drives)(nil)

// Drives accesses a map of drives in a bolt repository.
type Drives struct {
	db *bolt.DB
}

// Path returns the path of the drives.
func (ref Drives) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket}
}

// List returns the list of drives contained within the repository.
func (ref Drives) List() (ids []resource.ID, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(RootBucket))
		if root == nil {
			return nil
		}
		drives := root.Bucket([]byte(DriveBucket))
		if drives == nil {
			return nil
		}
		cursor := drives.Cursor()
		k, _ := cursor.First()
		for k != nil {
			ids = append(ids, resource.ID(k))
			k, _ = cursor.Next()
		}
		return nil
	})
	return ids, err
}

// Ref returns a drive reference.
func (ref Drives) Ref(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		db:    ref.db,
		drive: driveID,
	}
}
