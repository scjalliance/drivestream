package drivestream

import (
	"context"

	"github.com/scjalliance/drivestream/resource"
)

// A Collector is capable of collecting drivestream data.
type Collector interface {
	ChangeToken(ctx context.Context) (startToken string, err error)

	// Drive collects the current drive data, formatted in the same manner as
	// a change.
	Drive(ctx context.Context) (c resource.Change, err error)

	// Files collects a set of files into p, formatted in the same manner as
	// changes, up to len(p), starting from the file identified by token.
	//
	// If the provided token is empty it will start at the first file within
	// the drive. If len(p) is zero it will panic.
	//
	// The number of files collected are returned in n.
	//
	// If there more files to be collected in the current listing, nextToken
	// will be non-empty.
	Files(ctx context.Context, token string, p []resource.Change) (n int, nextToken string, err error)

	// Changes collects a set of changes into p, up to len(p), starting from
	// the change identified by token.
	//
	// If len(p) is zero it will panic.
	//
	// The number of changes collected are returned in n.
	//
	// If there more changes to be collected in the current set, nextToken
	// will be non-empty. If there are no more changes in the current set
	// then nextStartToken will hold the starting token for the next set.
	Changes(ctx context.Context, token string, p []resource.Change) (n int, nextToken string, nextStartToken string, err error)
}
