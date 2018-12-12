package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.TreeGroup = (*CommitTreeGroup)(nil)

// CommitTreeGroup is an unordered group of tree changes sharing a common
// parent.
type CommitTreeGroup struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
	parent resource.ID
}

// Parent returns the parent resource ID of the group.
func (ref CommitTreeGroup) Parent() resource.ID {
	return ref.parent
}

// Changes returns the set of changes contained in the group.
func (ref CommitTreeGroup) Changes() (changes []commit.TreeChange, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	tree := drv.Commits[ref.commit].Tree
	if tree == nil {
		return nil, commit.TreeGroupNotFound{Drive: ref.drive, Commit: ref.commit, Parent: ref.parent}
	}
	group := tree[ref.parent]
	if group == nil {
		return nil, commit.TreeGroupNotFound{Drive: ref.drive, Commit: ref.commit, Parent: ref.parent}
	}
	for child, removed := range group {
		changes = append(changes, commit.TreeChange{
			Parent:  ref.parent,
			Child:   child,
			Removed: removed,
		})
	}
	return changes, nil
}
