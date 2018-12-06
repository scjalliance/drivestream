package main

import (
	"io"

	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/boltrepo"
	"github.com/scjalliance/drivestream/memrepo"
	"github.com/scjalliance/drivestream/resource"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	boltType = "bolt"
	memType  = "in-memory"
)

// A creator is capable of creating drivestream repositories.
type creator func(resource.ID) drivestream.Repository

// DB holds database information for the application.
type DB struct {
	kind         string
	create       creator
	closer       io.Closer
	repositories map[resource.ID]drivestream.Repository
}

// NewDB returns an instance of the requested database.
func NewDB(app *kingpin.Application, dbType, path string) *DB {
	db := DB{
		repositories: make(map[resource.ID]drivestream.Repository),
	}

	switch dbType {
	case "bolt":
		boltDB, err := bolt.Open(path, 0600, nil)
		if err != nil {
			app.Errorf("failed to create or open bolt database: %v", err)
		}
		db.kind = "bolt"
		db.create = func(id resource.ID) drivestream.Repository {
			return boltrepo.New(boltDB, id)
		}
		db.closer = boltDB
	case "in-memory", "mem", "memory":
		db.kind = "in-memory"
		db.create = func(id resource.ID) drivestream.Repository {
			return memrepo.New(id)
		}
	default:
		app.Errorf("unrecognized database type: %s", dbType)
	}

	return &db
}

// Kind returns the kind of database db is.
func (db *DB) Kind() string {
	return db.kind
}

// Repository returns a drivestream repository instance for the given
// team drive ID, and reports whether the directory already existed.
func (db *DB) Repository(id resource.ID) (repo drivestream.Repository, existing bool) {
	repo, existing = db.repositories[id]
	if !existing {
		repo = db.create(id)
		db.repositories[id] = repo
	}
	return
}

// Close releases any resource consumed by the selected database
func (db *DB) Close() (err error) {
	if db.closer == nil {
		return nil
	}
	return db.closer.Close()
}
