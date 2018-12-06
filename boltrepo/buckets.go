package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

// collectionsBucket returns the collections bucket of the team drive.
func collectionsBucket(tx *bolt.Tx, teamDriveID resource.ID) *bolt.Bucket {
	root := tx.Bucket([]byte(RootBucket))
	if root == nil {
		return nil
	}
	anchor := root.Bucket([]byte(teamDriveID))
	if anchor == nil {
		return nil
	}
	return anchor.Bucket([]byte(CollectionBucket))
}

// collectionsBucket returns the collections bucket of the team drive.
func createCollectionsBucket(tx *bolt.Tx, teamDriveID resource.ID) (*bolt.Bucket, error) {
	root, err := tx.CreateBucketIfNotExists([]byte(RootBucket))
	if err != nil {
		return nil, err
	}
	anchor, err := root.CreateBucketIfNotExists([]byte(teamDriveID))
	if err != nil {
		return nil, err
	}
	return anchor.CreateBucketIfNotExists([]byte(CollectionBucket))
}

// collectionBucket returns the bucket of a particular collection.
func collectionBucket(tx *bolt.Tx, teamDriveID resource.ID, c collection.SeqNum) *bolt.Bucket {
	collections := collectionsBucket(tx, teamDriveID)
	if collections == nil {
		return nil
	}
	key := makeCollectionKey(c)
	return collections.Bucket(key[:])
}
