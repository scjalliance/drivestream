package commit

import "github.com/scjalliance/drivestream/resource"

// TreeGroup is an unordered group of tree changes sharing a common parent.
type TreeGroup interface {
	// Parent returns the parent resource ID of the group.
	Parent() resource.ID

	// Changes returns the set of changes contained in the group.
	Changes() (changes []TreeChange, err error)
}
