package commit

import "github.com/scjalliance/drivestream/resource"

// TreeMap is an unordered map of tree changes, grouped by parent.
type TreeMap interface {
	// Parents returns a list of parent IDs contained within the map.
	Parents() (parents []resource.ID, err error)

	// Group returns a reference to a group of changes sharing parent.
	Group(parent resource.ID) TreeGroup

	// Add adds the given tree changes to the map, grouped by parent.
	// If two or more changes conflict, the last change added takes
	// precedence.
	Add(changes ...TreeChange) error
}
