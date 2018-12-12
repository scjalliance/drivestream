package driveversion

import "github.com/scjalliance/drivestream/resource"

// A Sequence is an ordered series of drive versions.
type Sequence interface {
	// Next returns the next version number in the sequence.
	Next() (n resource.Version, err error)

	// Read reads drive data for a range of drive versions starting at the
	// given version number. Up to len(p) entries will be returned in p.
	// The number of entries is returned as n.
	Read(start resource.Version, p []resource.DriveData) (n int, err error)

	// Ref returns a drive version reference for the version number.
	Ref(v resource.Version) Reference
}
