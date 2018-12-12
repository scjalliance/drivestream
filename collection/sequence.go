package collection

// A Sequence is an ordered series of drivestream collections.
type Sequence interface {
	// Next returns the sequence number to use for the next collection.
	Next() (n SeqNum, err error)

	// Read reads collection data for a range of collections
	// starting at the given sequence number. Up to len(p) entries will
	// be returned in p. The number of entries is returned as n.
	Read(start SeqNum, p []Data) (n int, err error)

	// Ref returns a collection reference for the sequence number.
	Ref(seqNum SeqNum) Reference
}
