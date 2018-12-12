package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.StateSequence = (*CommitStates)(nil)

// CommitStates accesses a sequence of commit states in an in-memory
// repository.
type CommitStates struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
}

// Next returns the state number to use for the next state.
func (seq CommitStates) Next() (n commit.StateNum, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, nil
	}
	if seq.commit >= commit.SeqNum(len(drv.Commits)) {
		return 0, nil
	}
	return commit.StateNum(len(drv.Commits[seq.commit].States)), nil
}

// Read reads a subset of states from the sequence, starting at start.
// Up to len(p) states will be returned in p. The number of states
// returned is provided as n.
func (seq CommitStates) Read(start commit.StateNum, p []commit.State) (n int, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, commit.NotFound{Drive: seq.drive, Commit: seq.commit}
	}
	if seq.commit >= commit.SeqNum(len(drv.Commits)) {
		return 0, commit.NotFound{Drive: seq.drive, Commit: seq.commit}
	}
	length := commit.StateNum(len(drv.Commits[seq.commit].States))
	if start >= length {
		return 0, commit.StateNotFound{Drive: seq.drive, Commit: seq.commit, State: start}
	}
	for n < len(p) && start+commit.StateNum(n) < length {
		p[n] = drv.Commits[seq.commit].States[start+commit.StateNum(n)]
		n++
	}
	return n, nil
}

// Ref returns a commit state reference for the sequence number.
func (seq CommitStates) Ref(stateNum commit.StateNum) commit.StateReference {
	return CommitState{
		repo:   seq.repo,
		drive:  seq.drive,
		commit: seq.commit,
		state:  stateNum,
	}
}
