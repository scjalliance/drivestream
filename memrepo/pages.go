package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ page.Sequence = (*Pages)(nil)

// Pages accesses a sequence of pages in an in-memory repository.
type Pages struct {
	repo       *Repository
	drive      resource.ID
	collection collection.SeqNum
}

// Next returns the sequence number to use for the next page of the
// collection.
func (seq Pages) Next() (n page.SeqNum, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, nil
	}
	if seq.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, nil
	}
	return page.SeqNum(len(drv.Collections[seq.collection].Pages)), nil
}

// Read reads the requested pages from a collection.
func (seq Pages) Read(start page.SeqNum, p []page.Data) (n int, err error) {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return 0, collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	if seq.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	length := page.SeqNum(len(drv.Collections[seq.collection].Pages))
	if start >= length {
		return 0, collection.PageNotFound{Drive: seq.drive, Collection: seq.collection, Page: start}
	}
	for n < len(p) && start+page.SeqNum(n) < length {
		p[n] = drv.Collections[seq.collection].Pages[start+page.SeqNum(n)]
		n++
	}
	return n, nil
}

// Ref returns a page reference for the sequence number.
func (seq Pages) Ref(pageNum page.SeqNum) page.Reference {
	return Page{
		repo:       seq.repo,
		drive:      seq.drive,
		collection: seq.collection,
		page:       pageNum,
	}
}

// Clear removes all pages affiliated with a collection.
func (seq Pages) Clear() error {
	drv, ok := seq.repo.drives[seq.drive]
	if !ok {
		return collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	if seq.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.NotFound{Drive: seq.drive, Collection: seq.collection}
	}
	drv.Collections[seq.collection].Pages = nil
	return nil
}
