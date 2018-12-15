package boltrepo

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.FileMap = (*Files)(nil)

// Files accesses a map of files in a bolt repository.
type Files struct {
	db *bolt.DB
}

// List returns the list of files contained within the repository.
func (repo Files) List() (ids []resource.ID, err error) {
	err = repo.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(RootBucket))
		if root == nil {
			return nil
		}
		files := root.Bucket([]byte(FileBucket))
		if files == nil {
			return nil
		}
		cursor := files.Cursor()
		k, _ := cursor.First()
		for k != nil {
			ids = append(ids, resource.ID(k))
			k, _ = cursor.Next()
		}
		return nil
	})
	return ids, err
}

// Ref returns a file reference.
func (repo Files) Ref(id resource.ID) drivestream.FileReference {
	return File{
		db:   repo.db,
		file: id,
	}
}

// AddVersions adds file versions to the file map in bulk.
func (repo Files) AddVersions(fileVersions ...resource.File) error {
	// Perform the JSON encoding outside the transaction to minimize
	// time spent within it.
	payloads := make([][]byte, len(fileVersions))
	for i := range fileVersions {
		payload, err := json.Marshal(fileVersions[i].FileData)
		if err != nil {
			return err
		}
		payloads = append(payloads, payload)
	}
	return repo.db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(RootBucket))
		if err != nil {
			return err
		}
		files, err := root.CreateBucketIfNotExists([]byte(FileBucket))
		if err != nil {
			return err
		}
		for i := range fileVersions {
			file, err := files.CreateBucketIfNotExists([]byte(fileVersions[i].ID))
			if err != nil {
				return err
			}
			versions, err := file.CreateBucketIfNotExists([]byte(VersionBucket))
			if err != nil {
				return err
			}
			key := makeVersionKey(fileVersions[i].Version)
			if err := versions.Put(key[:], payloads[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

// AddViewData adds view data to the file map in bulk.
func (repo Files) AddViewData(entries ...fileview.Data) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		for _, entry := range entries {
			views, err := createFileViewsBucket(tx, entry.File)
			if err != nil {
				return err
			}

			view, err := views.CreateBucketIfNotExists([]byte(entry.Drive))
			if err != nil {
				return err
			}

			key := makeCommitKey(entry.Commit)
			value := makeVersionKey(entry.Version)
			err = view.Put(key[:], value[:])
			if err != nil {
				return err
			}
		}
		return nil
	})
}
