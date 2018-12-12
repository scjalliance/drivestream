package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// CommitEntry holds data for a commit.
type CommitEntry struct {
	Data   commit.Data
	States []commit.State
	Files  map[resource.ID]resource.Version
	Tree   map[resource.ID]map[resource.ID]bool
}

func newCommitEntry(data commit.Data) CommitEntry {
	return CommitEntry{
		Data: data,
	}
}
