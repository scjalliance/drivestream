package seqnum

// Number is a typed 64-bit sequence number.
type Number interface {
	Type() Type
	Value() int64
}
