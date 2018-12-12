package memrepo

import (
	"sort"

	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.DriveMap = (*Drives)(nil)

// Drives accesses a map of drives in an in-memory repository.
type Drives struct {
	repo *Repository
}

// List returns the list of drives contained within the repository.
func (driveMap Drives) List() (ids []resource.ID, err error) {
	for id := range driveMap.repo.drives {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}

// Ref returns a drive reference.
func (driveMap Drives) Ref(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		repo:  driveMap.repo,
		drive: driveID,
	}
}
