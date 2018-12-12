package commit

// A Sequence is an ordered series of drivestream commits.
type Sequence interface {
	// Next returns the sequence number to use for the next commit.
	Next() (n SeqNum, err error)

	// Read reads commit data for a range of commits
	// starting at the given sequence number. Up to len(p) entries will
	// be returned in p. The number of entries is returned as n.
	Read(start SeqNum, p []Data) (n int, err error)

	// Ref returns a collection reference for the sequence number.
	Ref(seqNum SeqNum) Reference
}
