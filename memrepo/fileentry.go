package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// FileEntry holds version history for a file.
type FileEntry struct {
	Versions map[resource.Version]resource.FileData
	Views    map[resource.ID]map[commit.SeqNum]resource.Version
}

func newFileEntry() FileEntry {
	return FileEntry{
		Versions: make(map[resource.Version]resource.FileData),
		Views:    make(map[resource.ID]map[commit.SeqNum]resource.Version),
	}
}
