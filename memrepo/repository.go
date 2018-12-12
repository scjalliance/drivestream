package memrepo

import (
	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/resource"
)

var _ drivestream.Repository = (*Repository)(nil)

// Repository is an in-memory implementation of a drive stream repository.
// It should be created by calling New.
type Repository struct {
	drives map[resource.ID]DriveEntry
	files  map[resource.ID]FileEntry
	//files       map[resource.ID]File
	//trees       map[resource.ID]Tree
	//content     map[filetree.Hash]filetree.Content
	//changeSets []drivestream.ChangeSet
}

// New returns a new in-memory drivestream repository.
func New() *Repository {
	return &Repository{
		drives: make(map[resource.ID]DriveEntry),
		files:  make(map[resource.ID]FileEntry),
		//trees:   make(map[resource.ID]Tree),
		//content: make(map[filetree.Hash]filetree.Content),
	}
}

// Type returns a string describing the type of the repository.
func (repo *Repository) Type() string {
	return "in-memory"
}

// Drives returns a drive map.
func (repo *Repository) Drives() drivestream.DriveMap {
	return Drives{repo: repo}
}

// Drive returns a drive reference.
func (repo *Repository) Drive(driveID resource.ID) drivestream.DriveReference {
	return Drive{
		repo:  repo,
		drive: driveID,
	}
}

// Files returns a file map.
func (repo *Repository) Files() drivestream.FileMap {
	return Files{repo: repo}
}

// File returns a file reference.
func (repo *Repository) File(fileID resource.ID) drivestream.FileReference {
	return File{
		repo: repo,
		file: fileID,
	}
}
