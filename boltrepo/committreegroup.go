package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.TreeGroup = (*CommitTreeGroup)(nil)

// CommitTreeGroup is an unordered group of tree changes sharing a common
// parent.
type CommitTreeGroup struct {
	db     *bolt.DB
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
	err = ref.db.View(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		tree := com.Bucket([]byte(TreeBucket))
		if tree == nil {
			return commit.TreeGroupNotFound{Drive: ref.drive, Commit: ref.commit, Parent: ref.parent}
		}
		group := tree.Bucket([]byte(ref.parent))
		if group == nil {
			return commit.TreeGroupNotFound{Drive: ref.drive, Commit: ref.commit, Parent: ref.parent}
		}
		cursor := group.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if len(v) != 1 {
				return commit.TreeGroupInvalid{Drive: ref.drive, Commit: ref.commit, Parent: ref.parent}
			}
			changes = append(changes, commit.TreeChange{
				Parent:  ref.parent,
				Child:   resource.ID(k),
				Removed: v[0] != 0,
			})
		}
		return nil
	})
	return changes, err
}
