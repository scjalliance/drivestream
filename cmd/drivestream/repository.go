package main

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/boltrepo"
	"github.com/scjalliance/drivestream/memrepo"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// NewRepository returns an instance of the requested database.
func NewRepository(app *kingpin.Application, dbType, path string) (drivestream.Repository, func() error) {
	switch dbType {
	case "bolt":
		boltDB, err := bolt.Open(path, 0600, nil)
		if err != nil {
			app.Errorf("failed to create or open bolt database: %v", err)
		}
		return boltrepo.New(boltDB), boltDB.Close
	case "in-memory", "mem", "memory":
		return memrepo.New(), func() error { return nil }
	default:
		app.Fatalf("unrecognized database type: %s", dbType)
		return nil, nil
	}
}
