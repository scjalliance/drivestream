package main

import (
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// driveData returns the most recent drive data for the repository.
func driveData(repo drivestream.Repository) (data resource.DriveData, ok bool) {
	cursor, err := collection.NewCursor(repo)
	if err != nil {
		return resource.DriveData{}, false
	}

	// FIXME: Check incremental changes in addition to drive pages
	for cursor.Last(); cursor.Valid(); cursor.Previous() {
		r, err := cursor.Reader()
		if err != nil {
			continue
		}
		data, err := r.Data()
		if err != nil {
			continue
		}
		if data.Type != collection.Full {
			continue
		}
		if r.NextPage() == 0 {
			continue
		}
		pg, err := r.Page(0)
		if err != nil || pg.Type != page.DriveList {
			continue
		}
		if len(pg.Changes) == 0 {
			continue
		}
		return pg.Changes[len(pg.Changes)-1].DriveData, true
	}

	return resource.DriveData{}, false
}
