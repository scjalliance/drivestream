package memrepo

import "github.com/scjalliance/drivestream/commit"

// Commit holds data for a commit.
type Commit struct {
	Data   commit.Data
	States []commit.State
}
