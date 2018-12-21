package boltrepo

import (
	"github.com/boltdb/bolt"
	"github.com/scjalliance/drivestream/binpath"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.TreeMap = (*CommitTree)(nil)

// CommitTree is a reference to a commit file map.
type CommitTree struct {
	db     *bolt.DB
	drive  resource.ID
	commit commit.SeqNum
}

// Path returns the path of the commit tree.
func (ref CommitTree) Path() binpath.Text {
	return binpath.Text{RootBucket, DriveBucket, string(ref.drive), TreeBucket, ref.commit.String()}
}

// Parents returns a list of parent IDs contained within the map.
func (ref CommitTree) Parents() (parents []resource.ID, err error) {
	err = ref.db.View(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		tree := com.Bucket([]byte(TreeBucket))
		if tree == nil {
			return nil
		}
		cursor := tree.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if v != nil {
				return commit.TreeGroupInvalid{Drive: ref.drive, Commit: ref.commit, Parent: resource.ID(k)} // All groups must be buckets
			}
			parents = append(parents, resource.ID(k))
		}
		return nil
	})
	return parents, err
}

// Group returns a reference to a group of changes sharing parent.
func (ref CommitTree) Group(parent resource.ID) commit.TreeGroup {
	return CommitTreeGroup{
		db:     ref.db,
		drive:  ref.drive,
		commit: ref.commit,
		parent: parent,
	}
}

// Add adds the given tree changes to the map, grouped by parent.
// If two or more changes conflict, the last change added takes
// precedence.
func (ref CommitTree) Add(changes ...commit.TreeChange) error {
	return ref.db.Update(func(tx *bolt.Tx) error {
		com := commitBucket(tx, ref.drive, ref.commit)
		if com == nil {
			return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
		}
		tree, err := com.CreateBucketIfNotExists([]byte(TreeBucket))
		if err != nil {
			return err
		}
		for _, change := range changes {
			group, err := tree.CreateBucketIfNotExists([]byte(change.Parent))
			if err != nil {
				return err
			}
			key := []byte(change.Child)
			value := makeBool(change.Removed)
			if err := group.Put(key, value[:]); err != nil {
				return err
			}
		}
		return nil
	})
}
