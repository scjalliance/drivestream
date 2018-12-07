package commit

// Repository is an interface capable of accessing drivestream
// commit data.
type Repository interface {
	// NextCommit returns the sequence number to use for the next
	// commit.
	NextCommit() (n SeqNum, err error)

	// Commits returns commit data for a range of commits
	// starting at the given sequence number. Up to len(p) entries will
	// be returned in p. The number of entries is returned as n.
	Commits(start SeqNum, p []Data) (n int, err error)

	// CreateCommit creates a new commit with the given sequence
	// number and data. If a commit already exists with the sequence
	// number an error will be returned.
	CreateCommit(c SeqNum, data Data) error

	// NextCommitState returns the state number to use for the next
	// state of the commit.
	NextCommitState(c SeqNum) (n StateNum, err error)

	// CommitStates returns a range of commit states for the given
	// commit, starting at the given state number. Up to len(p) states
	// will be returned in p. The number of states is returned as n.
	CommitStates(c SeqNum, start StateNum, p []State) (n int, err error)

	// CreateCommitState creates a new commit state with the given
	// state number and data. If a state already exists with the state
	// number an error will be returned.
	CreateCommitState(c SeqNum, stateNum StateNum, state State) error
}
