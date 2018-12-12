package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.Repository = (*Repository)(nil)

// Repository is a drive stream repository backed by a bolt database.
// It should be created by calling New.
type Repository struct {
	db *bolt.DB
}

// New returns a new drivestream bolt database for the team drive.
func New(db *bolt.DB) Repository {
	return Repository{
		db: db,
	}
}

// Type returns a string describing the type of the repository.
func (repo Repository) Type() string {
	return "bolt"
}

// Drives returns a drive map.
func (repo Repository) Drives() drivestream.DriveMap {
	return Drives{db: repo.db}
}

// Drive returns a drive reference.
func (repo Repository) Drive(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		db:    repo.db,
		drive: driveID,
	}
}

// Files returns a file map.
func (repo Repository) Files() drivestream.FileMap {
	return Files{db: repo.db}
}

// File returns a file reference.
func (repo Repository) File(fileID resource.ID) drivestream.FileReference {
	return File{
		db:   repo.db,
		file: fileID,
	}
}
