package drivestream

import "github.com/scjalliance/drivestream/resource"

// FileMap is a map of drivestream files.
type FileMap interface {
	// Ref returns a file reference.
	Ref(fileID resource.ID) FileReference

	// AddVersions adds one or more file versions to the file map.
	AddVersions(files ...resource.File) error
}
