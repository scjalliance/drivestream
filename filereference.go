package drivestream

import (
	"github.com/scjalliance/drivestream/fileversion"
	"github.com/scjalliance/drivestream/fileview"
	"github.com/scjalliance/drivestream/resource"
)

// FileReference is a reference to a drivestream file.
type FileReference interface {
	// FileID returns the resource ID of the file.
	FileID() resource.ID

	// Exists returns true if the file exists.
	Exists() (bool, error)

	// Versions returns the version map for the file.
	Versions() fileversion.Map

	// Version returns a file version reference. Equivalent to Versions().Ref(s).
	Version(v resource.Version) fileversion.Reference

	// AddVersion adds a version to the file.
	//AddVersion(v resource.Version, data resource.FileData) error

	// Views returns the view map for the file.
	Views() fileview.Map

	// View returns a view of the file for a particular drive.
	View(driveID resource.ID) fileview.Reference

	// Tree returns the file tree for the file.
	//Tree() filetree.Tree
}
