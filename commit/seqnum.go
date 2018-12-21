package commit

import "github.com/scjalliance/drivestream/seqnum"

// A SeqNum is a commit sequence number.
type SeqNum int64

// String returns a string representation of the sequence number.
func (number SeqNum) String() string {
	v := number.Base64()
	return string(v[:])
}

// Base64 returns a base64 representation of the sequence number.
func (number SeqNum) Base64() (k [12]byte) {
	return seqnum.Encode(number)
}

// Type returns the type code for the sequence number.
func (number SeqNum) Type() seqnum.Type {
	return seqnum.DriveCommit
}

// Value returns the sequence number as 64-bit integer.
func (number SeqNum) Value() int64 {
	return int64(number)
}
