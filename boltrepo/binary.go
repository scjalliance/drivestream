package boltrepo

import (
	"encoding/binary"

	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// makeCollectionKey returns an 8 byte big-endian binary representation of
// a collection sequence number.
func makeCollectionKey(seqNum collection.SeqNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(seqNum))
	return
}

// makeCollectionStateKey returns an 8 byte big-endian binary representation
// of a collection state number.
func makeCollectionStateKey(stateNum collection.StateNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(stateNum))
	return
}

// makePageKey returns an 8 byte big-endian binary representation of
// a page sequence number.
func makePageKey(seqNum page.SeqNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(seqNum))
	return
}

// makeCommitKey returns an 8 byte big-endian binary representation of
// a commit sequence number.
func makeCommitKey(seqNum commit.SeqNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(seqNum))
	return
}

// makeCommitStateKey returns an 8 byte big-endian binary representation
// of a commit state number.
func makeCommitStateKey(stateNum commit.StateNum) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(stateNum))
	return
}

// makeVersionKey returns an 8 byte big-endian binary representation
// of a version number.
func makeVersionKey(version resource.Version) (key [8]byte) {
	binary.BigEndian.PutUint64(key[:], uint64(version))
	return
}
