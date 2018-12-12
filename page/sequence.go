package page

// A Sequence is an ordered series of drivestream pages.
type Sequence interface {
	// Next returns the sequence number to use for the next page of the
	// collection.
	Next() (n SeqNum, err error)

	// Last returns a reference to the last page of the collection.
	//Last() (Reference, error)

	// Read reads the requested pages from a collection.
	Read(start SeqNum, p []Data) (n int, err error)

	// Ref returns a page reference for the sequence number.
	Ref(seqNum SeqNum) Reference

	// Clear removes all pages affiliated with a collection.
	Clear() error
}
