package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.DriveMap = (*Drives)(nil)

// Drives accesses a map of drives in a bolt repository.
type Drives struct {
	db *bolt.DB
}

// List returns the list of drives contained within the repository.
func (driveMap Drives) List() (ids []resource.ID, err error) {
	err = driveMap.db.View(func(tx *bolt.Tx) error {
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
func (driveMap Drives) Ref(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		db:    driveMap.db,
		drive: driveID,
	}
}
