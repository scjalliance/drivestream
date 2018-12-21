package collection

import "github.com/scjalliance/drivestream/seqnum"

// StateNum is a collection state number.
type StateNum int64

// String returns a string representation of the state number.
func (number StateNum) String() string {
	v := number.Base64()
	return string(v[:])
}

// Base64 returns a base64 representation of the state number.
func (number StateNum) Base64() (k [12]byte) {
	return seqnum.Encode(number)
}

// Type returns the type code for the state number.
func (number StateNum) Type() seqnum.Type {
	return seqnum.DriveCollectionState
}

// Value returns the state number as 64-bit integer.
func (number StateNum) Value() int64 {
	return int64(number)
}
