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
func (ref Drives) List() (ids []resource.ID, err error) {
	for id := range ref.repo.drives {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}

// Ref returns a drive reference.
func (ref Drives) Ref(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		repo:  ref.repo,
		drive: driveID,
	}
}
