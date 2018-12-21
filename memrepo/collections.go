package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/resource"
)

var _ collection.Sequence = (*Collections)(nil)

// Collections accesses a sequence of collections in an in-memory repository.
type Collections struct {
	repo  *Repository
	drive resource.ID
}

// Next returns the sequence number to use for the next collection.
func (ref Collections) Next() (n collection.SeqNum, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, nil
	}
	return collection.SeqNum(len(drv.Collections)), nil
}

// Read reads collection data for a range of collections
// starting at the given sequence number. Up to len(p) entries will
// be returned in p. The number of entries is returned as n.
func (ref Collections) Read(start collection.SeqNum, p []collection.Data) (n int, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, collection.NotFound{Drive: ref.drive, Collection: start}
	}
	length := collection.SeqNum(len(drv.Collections))
	if start >= length {
		return 0, collection.NotFound{Drive: ref.drive, Collection: start}
	}
	for n < len(p) && start+collection.SeqNum(n) < length {
		p[n] = drv.Collections[start+collection.SeqNum(n)].Data
		n++
	}
	return n, nil
}

// Ref returns a collection reference.
func (ref Collections) Ref(c collection.SeqNum) collection.Reference {
	return Collection{
		repo:       ref.repo,
		drive:      ref.drive,
		collection: c,
	}
}
