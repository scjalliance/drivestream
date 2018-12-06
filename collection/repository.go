package collection

import "github.com/scjalliance/drivestream/page"

// Repository is an interface capable of accessing drivestream
// collection data.
type Repository interface {
	// NextCollection returns the sequence number to use for the next
	// collection.
	NextCollection() (SeqNum, error)

	// Collections returns collection data for a range of collections
	// starting at the given sequence number. Up to len(p) entries will
	// be returned in p. The number of entries is returned as n.
	Collections(start SeqNum, p []Data) (n int, err error)

	// CreateCollection creates a new collection with the given sequence
	// number and data. If a collection already exists with the sequence
	// number an error will be returned.
	CreateCollection(SeqNum, Data) error

	// NextCollectionState returns the state number to use for the next
	// state of the collection.
	NextCollectionState(SeqNum) (StateNum, error)

	// CollectionStates returns a range of collection states for the given
	// collection, starting at the given state number. Up to len(p) states
	// will be returned in p. The number of states is returned as n.
	CollectionStates(col SeqNum, start StateNum, p []State) (n int, err error)

	// CreateCollectionState creates a new collection state with the given
	// state number and data. If a state already exists with the state
	// number an error will be returned.
	CreateCollectionState(SeqNum, StateNum, State) error

	// NextPage returns the sequence number to use for the next page of the
	// collection.
	NextPage(SeqNum) (page.SeqNum, error)

	// Pages returns the requested page from a collection.
	Pages(col SeqNum, start page.SeqNum, p []page.Data) (n int, err error)

	// CreatePage creates a new page within a collection.
	CreatePage(SeqNum, page.SeqNum, page.Data) error

	// ClearPages removes all pages affiliated with a collection.
	ClearPages(SeqNum) error
}
