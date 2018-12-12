package commit

import "github.com/scjalliance/drivestream/resource"

// TreeChange describes a tree change in a commit.
type TreeChange struct {
	Parent  resource.ID
	Child   resource.ID
	Removed bool
}
