package seqnum

import (
	"encoding/base64"
	"encoding/binary"
)

// Decode decodes a sequence number and its type from base64 url-encoded data.
func Decode(data []byte) (t Type, value int64) {
	if len(data) < 12 {
		return Invalid, 0
	}
	var raw [9]byte
	n, err := base64.URLEncoding.Decode(raw[:], data[:12])
	if n != len(raw) || err != nil {
		return Invalid, 0
	}
	return Type(raw[0]), int64(binary.BigEndian.Uint64(raw[1:]))
}
