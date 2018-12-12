package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.Reference = (*Commit)(nil)

// Commit is a drivestream commit reference for an in-memory repository.
type Commit struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
}

// Drive returns the drive ID of the commit.
func (ref Commit) Drive() resource.ID {
	return ref.drive
}

// SeqNum returns the sequence number of the commit.
func (ref Commit) SeqNum() commit.SeqNum {
	return ref.commit
}

// Exists returns true if the commit exists.
func (ref Commit) Exists() (exists bool, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return false, nil
	}
	return ref.commit < commit.SeqNum(len(drv.Commits)), nil
}

// Create creates a new commit with the given sequence number and data.
// If a commit already exists with the sequence number an error will be
// returned.
func (ref Commit) Create(data commit.Data) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		drv = newDriveEntry()
	}
	expected := commit.SeqNum(len(drv.Commits))
	if ref.commit != expected {
		return commit.OutOfOrder{Drive: ref.drive, Commit: ref.commit, Expected: expected}
	}
	drv.Commits = append(drv.Commits, newCommitEntry(data))
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns information about the commit.
func (ref Commit) Data() (data commit.Data, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return commit.Data{}, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return commit.Data{}, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	return drv.Commits[ref.commit].Data, nil
}

// States returns the state sequence for the commit.
func (ref Commit) States() commit.StateSequence {
	return CommitStates{
		repo:   ref.repo,
		drive:  ref.drive,
		commit: ref.commit,
	}
}

// State returns a state reference.
func (ref Commit) State(stateNum commit.StateNum) commit.StateReference {
	return ref.States().Ref(stateNum)
}

// Files returns the map of file changes for the commit.
func (ref Commit) Files() commit.FileMap {
	return CommitFiles{
		repo:   ref.repo,
		drive:  ref.drive,
		commit: ref.commit,
	}
}

// Tree returns the map of tree changes for the commit.
func (ref Commit) Tree() commit.TreeMap {
	return CommitTree{
		repo:   ref.repo,
		drive:  ref.drive,
		commit: ref.commit,
	}
}
