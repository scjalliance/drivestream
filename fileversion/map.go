package fileversion

import "github.com/scjalliance/drivestream/resource"

// A Map is a map of file versions.
type Map interface {
	// List returns a list of version numbers for the file.
	List() (v []resource.Version, err error)

	// Ref returns a file version reference for the version number.
	Ref(v resource.Version) Reference
}
