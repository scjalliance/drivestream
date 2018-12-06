package collection

import (
	"fmt"

	"github.com/scjalliance/drivestream/page"
)

// NotFound reports that a requested collection does not exist
// within the repository.
type NotFound struct {
	SeqNum SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream collection %d could not be found", e.SeqNum)
}

// OutOfOrder reports that a collection could not be created because its
// sequence number is not the next one in the series.
type OutOfOrder struct {
	SeqNum   SeqNum
	Expected SeqNum
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream collection %d could not be created because it would be out of order (expected %d)", e.SeqNum, e.Expected)
}

// InvalidData reports that a requested collection contains invalid
// or unparsable data.
type InvalidData struct {
	SeqNum SeqNum
}

// Error returns a string representation of the error.
func (e InvalidData) Error() string {
	return fmt.Sprintf("drivestream collection %d contains invalid data", e.SeqNum)
}

// StateNotFound reports that a requested state does not exist
// within the collection.
type StateNotFound struct {
	SeqNum   SeqNum
	StateNum StateNum
}

// Error returns a string representation of the error.
func (e StateNotFound) Error() string {
	return fmt.Sprintf("drivestream collection %d does not contain state %d", e.SeqNum, e.StateNum)
}

// StateOutOfOrder reports that a collection state could not be created
// because its state number is not the next one in the series.
type StateOutOfOrder struct {
	SeqNum   SeqNum
	StateNum StateNum
	Expected StateNum
}

// Error returns a string representation of the error.
func (e StateOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream collection %d state %d could not be created because it would be out of order (expected %d)", e.SeqNum, e.StateNum, e.Expected)
}

// InvalidState reports that a requested collection contains invalid
// or unparsable state.
type InvalidState struct {
	SeqNum   SeqNum
	StateNum StateNum
}

// Error returns a string representation of the error.
func (e InvalidState) Error() string {
	return fmt.Sprintf("drivestream collection %d contains invalid data for state %d", e.SeqNum, e.StateNum)
}

// PageNotFound reports that a requested page does not exist
// within the collection.
type PageNotFound struct {
	SeqNum  SeqNum
	PageNum page.SeqNum
}

// Error returns a string representation of the error.
func (e PageNotFound) Error() string {
	return fmt.Sprintf("drivestream collection %d does not contain page %d", e.SeqNum, e.PageNum)
}

// PageOutOfOrder reports that a collection page could not be created
// because its page number is not the next one in the series.
type PageOutOfOrder struct {
	SeqNum   SeqNum
	PageNum  page.SeqNum
	Expected page.SeqNum
}

// Error returns a string representation of the error.
func (e PageOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream collection %d page %d could not be created because it would be out of order (expected %d)", e.SeqNum, e.PageNum, e.Expected)
}

// InvalidPage reports that a requested collection contains invalid
// or unparsable page data.
type InvalidPage struct {
	SeqNum  SeqNum
	PageNum page.SeqNum
}

// Error returns a string representation of the error.
func (e InvalidPage) Error() string {
	return fmt.Sprintf("drivestream collection %d contains invalid data for state %d", e.SeqNum, e.PageNum)
}
