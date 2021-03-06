package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.Sequence = (*Commits)(nil)

// Commits accesses a sequence of commits in an in-memory repository.
type Commits struct {
	repo  *Repository
	drive resource.ID
}

// Next returns the sequence number to use for the next commit.
func (ref Commits) Next() (n commit.SeqNum, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, nil
	}
	return commit.SeqNum(len(drv.Commits)), nil
}

// Read reads commit data for a range of commits
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (ref Commits) Read(start commit.SeqNum, p []commit.Data) (n int, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, commit.NotFound{Drive: ref.drive, Commit: start}
	}
	length := commit.SeqNum(len(drv.Commits))
	if start >= length {
		return 0, commit.NotFound{Drive: ref.drive, Commit: start}
	}
	for n < len(p) && start+commit.SeqNum(n) < length {
		p[n] = drv.Commits[start+commit.SeqNum(n)].Data
		n++
	}
	return n, nil
}

// Ref returns a commit reference.
func (ref Commits) Ref(c commit.SeqNum) commit.Reference {
	return Commit{
		repo:   ref.repo,
		drive:  ref.drive,
		commit: c,
	}
}
