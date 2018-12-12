package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
)

// CollectionEntry holds data for a collection.
type CollectionEntry struct {
	Data   collection.Data
	States []collection.State
	Pages  []page.Data
}

func newCollectionEntry(data collection.Data) CollectionEntry {
	return CollectionEntry{
		Data: data,
	}
}
