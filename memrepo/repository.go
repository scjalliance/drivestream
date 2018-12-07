package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// Repository is an in-memory implementation of a drive stream repository.
// It should be created by calling New.
type Repository struct {
	drive       Drive
	collections []Collection
	commits     []Commit
	//files       map[resource.ID]File
	//trees       map[resource.ID]Tree
	//content     map[filetree.Hash]filetree.Content
	//changeSets []drivestream.ChangeSet
}

// New returns a new in-memory drivestream database for the team drive.
func New(teamDriveID resource.ID) *Repository {
	return &Repository{
		drive: Drive{
			ID: teamDriveID,
		},
		//files: make(map[resource.ID]File),
		//trees:   make(map[resource.ID]Tree),
		//content: make(map[filetree.Hash]filetree.Content),
	}
}

// DriveID returns the team drive ID that the drive stream is for.
func (repo *Repository) DriveID() resource.ID {
	return repo.drive.ID
}

// NextCollection returns the sequence number to use for the next
// collection.
func (repo *Repository) NextCollection() (collection.SeqNum, error) {
	return collection.SeqNum(len(repo.collections)), nil
}

// Collections returns collection data for a range of collections
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (repo *Repository) Collections(start collection.SeqNum, p []collection.Data) (n int, err error) {
	if int(start) >= len(repo.collections) {
		return 0, collection.NotFound{SeqNum: start}
	}
	for n < len(p) && n+int(start) < len(repo.collections) {
		p[n] = repo.collections[n+int(start)].Data
		n++
	}
	return n, nil
}

// CreateCollection creates a new collection with the given sequence
// number and data. If a collection already exists with the sequence
// number an error will be returned.
func (repo *Repository) CreateCollection(c collection.SeqNum, data collection.Data) error {
	if int(c) != len(repo.collections) {
		return collection.OutOfOrder{SeqNum: c}
	}
	repo.collections = append(repo.collections, Collection{Data: data})
	return nil
}

// NextCollectionState returns the state number to use for the next
// state of the collection.
func (repo *Repository) NextCollectionState(c collection.SeqNum) (collection.StateNum, error) {
	if int(c) >= len(repo.collections) {
		return 0, collection.NotFound{SeqNum: c}
	}
	return collection.StateNum(len(repo.collections[c].States)), nil
}

// CollectionStates returns a range of collection states for the given
// collection, starting at the given state number. Up to len(p) states
// will be returned in p. The number of states is returned as n.
func (repo *Repository) CollectionStates(c collection.SeqNum, start collection.StateNum, p []collection.State) (n int, err error) {
	if int(c) >= len(repo.collections) {
		return 0, collection.NotFound{SeqNum: c}
	}
	col := &repo.collections[c]

	if int(start) >= len(col.States) {
		return 0, collection.StateNotFound{SeqNum: c, StateNum: start}
	}

	return copy(p, col.States[start:]), nil
}

// CreateCollectionState creates a new collection state with the given
// state number and data. If a state already exists with the state
// number an error will be returned.
func (repo *Repository) CreateCollectionState(c collection.SeqNum, stateNum collection.StateNum, state collection.State) error {
	if int(c) >= len(repo.collections) {
		return collection.NotFound{SeqNum: c}
	}
	if int(stateNum) != len(repo.collections[c].States) {
		return collection.StateOutOfOrder{SeqNum: c, StateNum: stateNum}
	}
	repo.collections[c].States = append(repo.collections[c].States, state)
	return nil
}

// Pages returns the requested pages from a collection.
func (repo *Repository) Pages(c collection.SeqNum, start page.SeqNum, p []page.Data) (n int, err error) {
	if int(c) >= len(repo.collections) {
		return 0, collection.NotFound{SeqNum: c}
	}
	col := &repo.collections[c]

	if int(start) >= len(col.Pages) {
		return 0, collection.PageNotFound{SeqNum: c, PageNum: start}
	}

	return copy(p, col.Pages[start:]), nil
}

// NextPage returns the sequence number to use for the next page of the
// collection.
func (repo *Repository) NextPage(c collection.SeqNum) (page.SeqNum, error) {
	if int(c) >= len(repo.collections) {
		return 0, collection.NotFound{SeqNum: c}
	}
	return page.SeqNum(len(repo.collections[c].Pages)), nil
}

// CreatePage creates a new page within a collection.
func (repo *Repository) CreatePage(c collection.SeqNum, pageNum page.SeqNum, data page.Data) error {
	if int(c) >= len(repo.collections) {
		return collection.NotFound{SeqNum: c}
	}
	if int(pageNum) != len(repo.collections[c].Pages) {
		return collection.PageOutOfOrder{SeqNum: c, PageNum: pageNum}
	}
	repo.collections[c].Pages = append(repo.collections[c].Pages, data)
	return nil
}

// ClearPages removes pages affiliated with a collection.
func (repo *Repository) ClearPages(c collection.SeqNum) error {
	if int(c) >= len(repo.collections) {
		return collection.NotFound{SeqNum: c}
	}
	repo.collections[c].Pages = nil
	return nil
}

// NextCommit returns the sequence number to use for the next
// commit.
func (repo *Repository) NextCommit() (n commit.SeqNum, err error) {
	return commit.SeqNum(len(repo.commits)), nil
}

// Commits returns commit data for a range of commits
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (repo *Repository) Commits(start commit.SeqNum, p []commit.Data) (n int, err error) {
	if int(start) >= len(repo.commits) {
		return 0, commit.NotFound{SeqNum: start}
	}
	for n < len(p) && n+int(start) < len(repo.commits) {
		p[n] = repo.commits[n+int(start)].Data
		n++
	}
	return n, nil
}

// CreateCommit creates a new commit with the given sequence
// number and data. If a commit already exists with the sequence
// number an error will be returned.
func (repo *Repository) CreateCommit(c commit.SeqNum, data commit.Data) error {
	if int(c) != len(repo.commits) {
		return commit.OutOfOrder{SeqNum: c}
	}
	repo.commits = append(repo.commits, Commit{Data: data})
	return nil
}

// NextCommitState returns the state number to use for the next
// state of the commit.
func (repo *Repository) NextCommitState(c commit.SeqNum) (n commit.StateNum, err error) {
	if int(c) >= len(repo.commits) {
		return 0, commit.NotFound{SeqNum: c}
	}
	return commit.StateNum(len(repo.commits[c].States)), nil
}

// CommitStates returns a range of commit states for the given
// commit, starting at the given state number. Up to len(p) states
// will be returned in p. The number of states is returned as n.
func (repo *Repository) CommitStates(c commit.SeqNum, start commit.StateNum, p []commit.State) (n int, err error) {
	if int(c) >= len(repo.commits) {
		return 0, commit.NotFound{SeqNum: c}
	}
	col := &repo.commits[c]

	if int(start) >= len(col.States) {
		return 0, commit.StateNotFound{SeqNum: c, StateNum: start}
	}

	return copy(p, col.States[start:]), nil
}

// CreateCommitState creates a new commit state with the given
// state number and data. If a state already exists with the state
// number an error will be returned.
func (repo *Repository) CreateCommitState(c commit.SeqNum, stateNum commit.StateNum, state commit.State) error {
	if int(c) >= len(repo.commits) {
		return commit.NotFound{SeqNum: c}
	}
	if int(stateNum) != len(repo.commits[c].States) {
		return commit.StateOutOfOrder{SeqNum: c, StateNum: stateNum}
	}
	repo.commits[c].States = append(repo.commits[c].States, state)
	return nil
}
