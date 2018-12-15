package fileview

import "github.com/scjalliance/drivestream/resource"

// A Map is a map of file views.
type Map interface {
	// List returns a list of drives with a view of the file.
	List() (drives []resource.ID, err error)

	// Ref returns a view of the file for a particular drive.
	Ref(driveID resource.ID) Reference
}
