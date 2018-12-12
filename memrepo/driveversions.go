package memrepo

import (
	"github.com/scjalliance/drivestream/driveversion"
	"github.com/scjalliance/drivestream/resource"
)

var _ driveversion.Sequence = (*DriveVersions)(nil)

// DriveVersions accesses a sequence of drive versions in an in-memory repository.
type DriveVersions struct {
	repo  *Repository
	drive resource.ID
}

// Next returns the next version number in the sequence.
func (seq DriveVersions) Next() (n resource.Version, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, nil
	}
	return resource.Version(len(drv.Versions)), nil
}

// Read reads drive data for a range of drive versions starting at the
// given version number. Up to len(p) entries will be returned in p.
// The number of entries is returned as n.
func (seq DriveVersions) Read(start resource.Version, p []resource.DriveData) (n int, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, driveversion.NotFound{Drive: seq.drive, Version: start}
	}
	length := resource.Version(len(drv.Versions))
	if start >= length {
		return 0, driveversion.NotFound{Drive: seq.drive, Version: start}
	}
	for n < len(p) && start+resource.Version(n) < length {
		p[n] = drv.Versions[start+resource.Version(n)]
		n++
	}
	return n, nil
}

// Ref returns a drive version reference for the version number.
func (seq DriveVersions) Ref(v resource.Version) driveversion.Reference {
	return DriveVersion{
		repo:    seq.repo,
		drive:   seq.drive,
		version: v,
	}
}
