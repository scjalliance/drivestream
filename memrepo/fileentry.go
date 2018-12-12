package memrepo

import "github.com/scjalliance/drivestream/resource"

// FileEntry holds version history for a file.
type FileEntry struct {
	Versions map[resource.Version]resource.FileData
}

func newFileEntry() FileEntry {
	return FileEntry{
		Versions: make(map[resource.Version]resource.FileData),
	}
}
