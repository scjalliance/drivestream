package memrepo

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

var _ commit.FileMap = (*CommitFiles)(nil)

// CommitFiles is a reference to a commit file map.
type CommitFiles struct {
	repo   *Repository
	drive  resource.ID
	commit commit.SeqNum
}

// Read returns the set of file changes for the commit, in unspecified
// order.
func (ref CommitFiles) Read() (changes []commit.FileChange, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return nil, commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	files := drv.Commits[ref.commit].Files
	if files == nil {
		return nil, nil
	}
	for file, version := range files {
		changes = append(changes, commit.FileChange{
			File:    file,
			Version: version,
		})
	}
	return changes, nil
}

// Add adds the given file changes to the map.
// If two or more changes conflict, the last change added takes
// precedence.
func (ref CommitFiles) Add(changes ...commit.FileChange) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	if ref.commit >= commit.SeqNum(len(drv.Commits)) {
		return commit.NotFound{Drive: ref.drive, Commit: ref.commit}
	}
	files := drv.Commits[ref.commit].Files
	if files == nil {
		files = make(map[resource.ID]resource.Version)
		drv.Commits[ref.commit].Files = files
	}
	for _, change := range changes {
		files[change.File] = change.Version
	}
	return nil
}
