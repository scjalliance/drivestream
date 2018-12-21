package collection

import "github.com/scjalliance/drivestream/seqnum"

// SeqNum is a collection sequence number.
type SeqNum int64

// String returns a string representation of the sequence number.
func (number SeqNum) String() string {
	//return "col" + hex.EncodeToString(binary.BigEndian.)
	v := number.Base64()
	return string(v[:])
}

// Base64 returns a base64 representation of the sequence number.
func (number SeqNum) Base64() (k [12]byte) {
	return seqnum.Encode(number)
}

// Type returns the type code for the sequence number.
func (number SeqNum) Type() seqnum.Type {
	return seqnum.DriveCollection
}

// Value returns the sequence number as 64-bit integer.
func (number SeqNum) Value() int64 {
	return int64(number)
}
