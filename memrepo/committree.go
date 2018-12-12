package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.TreeMap = (*CommitTree)(nil)

// CommitTree is a reference to a commit file map.
type CommitTree struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
}

// Parents returns a list of parent IDs contained within the map.
func (ref CommitTree) Parents() (parents []resource.ID, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	tree := drv.Commits[ref.commit].Tree
	if tree == nil {
		return nil, nil
	}
	for parent := range tree {
		parents = append(parents, parent)
	}
	return parents, nil
}

// Group returns a reference to a group of changes sharing parent.
func (ref CommitTree) Group(parent resource.ID) commit.TreeGroup {
	return CommitTreeGroup{
		repo:   ref.repo,
		drive:  ref.drive,
		commit: ref.commit,
		parent: parent,
	}
}

// Add adds the given tree changes to the map, grouped by parent.
// If two or more changes conflict, the last change added takes
// precedence.
func (ref CommitTree) Add(changes ...commit.TreeChange) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	tree := drv.Commits[ref.commit].Tree
	if tree == nil {
		tree = make(map[resource.ID]map[resource.ID]bool)
		drv.Commits[ref.commit].Tree = tree
	}
	for _, change := range changes {
		group := tree[change.Parent]
		if group == nil {
			group = make(map[resource.ID]bool)
			tree[change.Parent] = group
		}
		group[change.Child] = change.Removed
	}
	return nil
}
