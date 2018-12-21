package seqnum

// Type is a type of sequence.
type Type byte

// Sequence types.
const (
	Invalid              Type = 0x00
	Resource             Type = 0x10
	ResourceVersion      Type = 0x11
	DriveCollection      Type = 0x20
	DriveCollectionState Type = 0x21
	DriveCollectionPage  Type = 0x22
	DriveCommit          Type = 0x40
	DriveCommitState     Type = 0x41
)
