package memrepo

import (
	"github.com/scjalliance/drivestream/resource"
)

// File holds version and commit history for a file.
type File struct {
	ID       resource.ID
	Versions map[resource.Version]resource.FileData
	//Commits  map[commit.SeqNum]resource.Version
	//Tree     map[commit.SeqNum]filetree.Hash
}
