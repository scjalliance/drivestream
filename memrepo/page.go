package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

var _ page.Reference = (*Page)(nil)

// Page is a drivestream page reference for a bolt repository.
type Page struct {
	repo       *Repository
	drive      resource.ID
	collection collection.SeqNum
	page       page.SeqNum
}

// SeqNum returns the sequence number of the page.
func (ref Page) SeqNum() page.SeqNum {
	return ref.page
}

// Create creates the page with the given data.
func (ref Page) Create(data page.Data) error {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	expected := page.SeqNum(len(drv.Collections[ref.collection].Pages))
	if ref.page != expected {
		return collection.PageOutOfOrder{Drive: ref.drive, Collection: ref.collection, Page: ref.page, Expected: expected}
	}
	drv.Collections[ref.collection].Pages = append(drv.Collections[ref.collection].Pages, data)
	ref.repo.drives[ref.drive] = drv
	return nil
}

// Data returns the page data.
func (ref Page) Data() (data page.Data, err error) {
	drv, ok := ref.repo.drives[ref.drive]
	if !ok {
		return page.Data{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.collection >= collection.SeqNum(len(drv.Collections)) {
		return page.Data{}, collection.NotFound{Drive: ref.drive, Collection: ref.collection}
	}
	if ref.page >= page.SeqNum(len(drv.Collections[ref.collection].Pages)) {
		return page.Data{}, collection.PageNotFound{Drive: ref.drive, Collection: ref.collection, Page: ref.page}
	}
	return drv.Collections[ref.collection].Pages[ref.page], nil
}
