package seqnum

import (
	"encoding/base64"
	"encoding/binary"
)

// Encode encodes n as base64 url-encoded data.
func Encode(n Number) (v [12]byte) {
	var data [9]byte
	data[0] = byte(n.Type())
	binary.BigEndian.PutUint64(data[1:], uint64(n.Value()))
	base64.URLEncoding.Encode(v[:], data[:])
	return
}
