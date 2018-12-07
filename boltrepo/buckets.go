package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// driveBucket returns the bucket of the drive.
func driveBucket(tx *bolt.Tx, teamDriveID resource.ID) *bolt.Bucket {
	root := tx.Bucket([]byte(RootBucket))
	if root == nil {
		return nil
	}
	drives := root.Bucket([]byte(DriveBucket))
	if drives == nil {
		return nil
	}
	return drives.Bucket([]byte(teamDriveID))
}

// createDriveBucket creates a bucket for the drive.
func createDriveBucket(tx *bolt.Tx, teamDriveID resource.ID) (*bolt.Bucket, error) {
	root, err := tx.CreateBucketIfNotExists([]byte(RootBucket))
	if err != nil {
		return nil, err
	}
	drives, err := root.CreateBucketIfNotExists([]byte(DriveBucket))
	if err != nil {
		return nil, err
	}
	return drives.CreateBucketIfNotExists([]byte(teamDriveID))
}

// collectionsBucket returns the collections bucket of the drive.
func collectionsBucket(tx *bolt.Tx, teamDriveID resource.ID) *bolt.Bucket {
	drv := driveBucket(tx, teamDriveID)
	if drv == nil {
		return nil
	}
	return drv.Bucket([]byte(CollectionBucket))
}

// createCollectionsBucket creates the collections bucket for the drive.
func createCollectionsBucket(tx *bolt.Tx, teamDriveID resource.ID) (*bolt.Bucket, error) {
	drv, err := createDriveBucket(tx, teamDriveID)
	if err != nil {
		return nil, err
	}
	return drv.CreateBucketIfNotExists([]byte(CollectionBucket))
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

// commitsBucket returns the commits bucket of the drive.
func commitsBucket(tx *bolt.Tx, teamDriveID resource.ID) *bolt.Bucket {
	drv := driveBucket(tx, teamDriveID)
	if drv == nil {
		return nil
	}
	return drv.Bucket([]byte(CommitBucket))
}

// commitsBucket returns the commits bucket of the drive.
func createCommitsBucket(tx *bolt.Tx, teamDriveID resource.ID) (*bolt.Bucket, error) {
	drv, err := createDriveBucket(tx, teamDriveID)
	if err != nil {
		return nil, err
	}
	return drv.CreateBucketIfNotExists([]byte(CommitBucket))
}

// commitBucket returns the bucket of a particular commit.
func commitBucket(tx *bolt.Tx, teamDriveID resource.ID, c commit.SeqNum) *bolt.Bucket {
	commits := commitsBucket(tx, teamDriveID)
	if commits == nil {
		return nil
	}
	key := makeCommitKey(c)
	return commits.Bucket(key[:])
}
