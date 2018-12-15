package drivestream

import (
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

// FileMap is a map of drivestream files.
type FileMap interface {
	// Ref returns a file reference.
	Ref(fileID resource.ID) FileReference

	// AddVersions adds file versions to the file map in bulk.
	AddVersions(files ...resource.File) error

	// AddViewData adds view data to the file map in bulk.
	AddViewData(entries ...fileview.Data) error
}
