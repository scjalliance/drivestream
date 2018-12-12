package commit

import "github.com/scjalliance/drivestream/resource"

// Reference is a drivestream commit reference.
type Reference interface {
	// Drive returns the drive ID of the commit.
	Drive() resource.ID

	// SeqNum returns the sequence number of the commit.
	SeqNum() SeqNum

	// Exists returns true if the commit exists.
	Exists() (bool, error)

	// Create creates the commit with the given data. If a commit already
	// exists with the sequence number an error will be returned.
	Create(data Data) error

	// Data returns information about the commit.
	Data() (Data, error)

	// States returns the state sequence for the commit.
	States() StateSequence

	// State returns a state reference. Equivalent to States().Ref(s).
	State(s StateNum) StateReference

	// Files returns the map of file changes for the commit.
	Files() FileMap

	// Tree returns the map of tree changes for the commit.
	Tree() TreeMap
}
