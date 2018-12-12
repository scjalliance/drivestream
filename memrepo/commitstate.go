package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.StateReference = (*CommitState)(nil)

// CommitState is a reference to a commit state.
type CommitState struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
	state  commit.StateNum
}

// StateNum returns the state number of the reference.
func (ref CommitState) StateNum() commit.StateNum {
	return ref.state
}

// Create creates the commit state with the given data. If a state already
// exists with the state number an error will be returned.
func (ref CommitState) Create(data commit.State) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	expected := commit.StateNum(len(drv.Commits[ref.commit].States))
	if ref.state != expected {
		return commit.StateOutOfOrder{Drive: ref.drive, Commit: ref.commit, State: ref.state, Expected: expected}
	}
	drv.Commits[ref.commit].States = append(drv.Commits[ref.commit].States, data)
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns the commit state data.
func (ref CommitState) Data() (data commit.State, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return commit.State{}, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return commit.State{}, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.state >= commit.StateNum(len(drv.Commits[ref.commit].States)) {
		return commit.State{}, commit.StateNotFound{Drive: ref.drive, Commit: ref.commit, State: ref.state}
	}
	return drv.Commits[ref.commit].States[ref.state], nil
}
