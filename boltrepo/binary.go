package boltrepo

import (
	"encoding/binary"

	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
)

// makeCollectionKey returns an 8 byte big-endian binary representation of
// a collection sequence number.
func makeCollectionKey(seqNum collection.SeqNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(seqNum))
	return
}

// makeStateKey returns an 8 byte big-endian binary representation of
// a collection state number.
func makeStateKey(stateNum collection.StateNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(stateNum))
	return
}

// makePageKey returns an 8 byte big-endian binary representation of
// a page sequence number.
func makePageKey(seqNum page.SeqNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(seqNum))
	return
}
