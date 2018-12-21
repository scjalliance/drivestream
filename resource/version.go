package resource

import (
	"github.com/scjalliance/drivestream/seqnum"
)

// Version is a file or drive version number.
type Version int64

// String returns a string representation of the version number.
func (number Version) String() string {
	v := number.Base64()
	return string(v[:])
}

// Base64 returns a base64 representation of the version number.
func (number Version) Base64() (k [12]byte) {
	return seqnum.Encode(number)
}

// Type returns the type code for the version number.
func (number Version) Type() seqnum.Type {
	return seqnum.ResourceVersion
}

// Value returns the version number as 64-bit integer.
func (number Version) Value() int64 {
	return int64(number)
}
