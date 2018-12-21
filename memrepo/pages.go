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
func (ref Pages) Next() (n page.SeqNum, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, nil
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, nil
	}
	return page.SeqNum(len(drv.Collections[ref.collection].Pages)), nil
}

// Read reads the requested pages from a collection.
func (ref Pages) Read(start page.SeqNum, p []page.Data) (n int, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return 0, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return 0, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	length := page.SeqNum(len(drv.Collections[ref.collection].Pages))
	if start >= length {
		return 0, collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: start}
	}
	for n < len(p) && start+page.SeqNum(n) < length {
		p[n] = drv.Collections[ref.collection].Pages[start+page.SeqNum(n)]
		n++
	}
	return n, nil
}

// Ref returns a page reference for the sequence number.
func (ref Pages) Ref(pageNum page.SeqNum) page.Reference {
	return Page{
		repo:       ref.repo,
		drive:      ref.drive,
		collection: ref.collection,
		page:       pageNum,
	}
}

// Clear removes all pages affiliated with a collection.
func (ref Pages) Clear() error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	drv.Collections[ref.collection].Pages = nil
	return nil
}
