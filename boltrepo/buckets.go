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

// driveVersionsBucket returns the versions bucket of the drive.
func driveVersionsBucket(tx *bolt.Tx, driveID resource.ID) *bolt.Bucket {
	drv := driveBucket(tx, driveID)
	if drv == nil {
		return nil
	}
	return drv.Bucket([]byte(VersionBucket))
}

// createDriveVersionsBucket creates the versions bucket for the drive.
func createDriveVersionsBucket(tx *bolt.Tx, driveID resource.ID) (*bolt.Bucket, error) {
	drv, err := createDriveBucket(tx, driveID)
	if err != nil {
		return nil, err
	}
	return drv.CreateBucketIfNotExists([]byte(VersionBucket))
}

// fileBucket returns the bucket of the file.
func fileBucket(tx *bolt.Tx, fileID resource.ID) *bolt.Bucket {
	root := tx.Bucket([]byte(RootBucket))
	if root == nil {
		return nil
	}
	files := root.Bucket([]byte(FileBucket))
	if files == nil {
		return nil
	}
	return files.Bucket([]byte(fileID))
}

// createFileBucket creates a bucket for the file.
func createFileBucket(tx *bolt.Tx, fileID resource.ID) (*bolt.Bucket, error) {
	root, err := tx.CreateBucketIfNotExists([]byte(RootBucket))
	if err != nil {
		return nil, err
	}
	files, err := root.CreateBucketIfNotExists([]byte(FileBucket))
	if err != nil {
		return nil, err
	}
	return files.CreateBucketIfNotExists([]byte(fileID))
}

// fileVersionsBucket returns the versions bucket of the file.
func fileVersionsBucket(tx *bolt.Tx, fileID resource.ID) *bolt.Bucket {
	file := fileBucket(tx, fileID)
	if file == nil {
		return nil
	}
	return file.Bucket([]byte(VersionBucket))
}

// createFileVersionsBucket creates the versions bucket for the file.
func createFileVersionsBucket(tx *bolt.Tx, fileID resource.ID) (*bolt.Bucket, error) {
	file, err := createFileBucket(tx, fileID)
	if err != nil {
		return nil, err
	}
	return file.CreateBucketIfNotExists([]byte(VersionBucket))
}
