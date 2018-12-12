package collection

import (
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// Reference is a drivestream collection reference.
type Reference interface {
	// Drive returns the drive ID of the collection.
	Drive() resource.ID

	// SeqNum returns the sequence number of the collection.
	SeqNum() SeqNum

	// Exists returns true if the collection exists.
	Exists() (bool, error)

	// Create creates a new collection with the given sequence number and data.
	// If a collection already exists with the sequence number an error will be
	// returned.
	Create(data Data) error

	// Data returns information about the collection.
	Data() (Data, error)

	// States returns the state sequence for the collection.
	States() StateSequence

	// State returns a state reference. Equivalent to States().Ref(s).
	State(s StateNum) StateReference

	// Pages returns the page sequence for the collection.
	Pages() page.Sequence

	// Page returns a page reference. Equivalent to Pages().Ref(s).
	Page(s page.SeqNum) page.Reference
}
