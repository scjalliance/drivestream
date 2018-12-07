package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/resource"
)

// Enumerate returns the list of team drives contained within db.
func Enumerate(db *bolt.DB) (ids []resource.ID, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(RootBucket))
		if root == nil {
			return nil
		}
		cursor := root.Cursor()
		k, _ := cursor.First()
		for k != nil {
			ids = append(ids, resource.ID(k))
			k, _ = cursor.Next()
		}
		return nil
	})
	return ids, err
}
