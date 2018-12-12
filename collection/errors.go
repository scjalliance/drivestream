package collection

import (
	"fmt"

	"github.com/scjalliance/drivestream/page"
	"github.com/scjalliance/drivestream/resource"
)

// NotFound reports that a collection could not be found
// within the repository.
type NotFound struct {
	Drive      resource.ID
	Collection SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d could not be found", e.Drive, e.Collection)
}

// OutOfOrder reports that a collection could not be created
// because its sequence number is not the next one in the series.
type OutOfOrder struct {
	Drive      resource.ID
	Collection SeqNum
	Expected   SeqNum
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d could not be created because it would be out of order (expected %d)", e.Drive, e.Collection, e.Expected)
}

// DataInvalid reports that collection contains invalid or
// unparsable data.
type DataInvalid struct {
	Drive      resource.ID
	Collection SeqNum
}

// Error returns a string representation of the error.
func (e DataInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d contains invalid data", e.Drive, e.Collection)
}

// StateNotFound reports that a requested state does not exist
// within the collection.
type StateNotFound struct {
	Drive      resource.ID
	Collection SeqNum
	State      StateNum
}

// Error returns a string representation of the error.
func (e StateNotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d does not contain state %d", e.Drive, e.Collection, e.State)
}

// StateOutOfOrder reports that a collection state could not be created
// because its state number is not the next one in the series.
type StateOutOfOrder struct {
	Drive      resource.ID
	Collection SeqNum
	State      StateNum
	Expected   StateNum
}

// Error returns a string representation of the error.
func (e StateOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d state %d could not be created because it would be out of order (expected %d)", e.Drive, e.Collection, e.State, e.Expected)
}

// StateInvalid reports that a requested collection contains invalid
// or unparsable state.
type StateInvalid struct {
	Drive      resource.ID
	Collection SeqNum
	State      StateNum
}

// Error returns a string representation of the error.
func (e StateInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d contains invalid data for state %d", e.Drive, e.Collection, e.State)
}

// StatesTruncated reports that a request for collection states
// returned a shorter list than expected. This is an indication that the
// repository is providing an inconsistent view of its data.
type StatesTruncated struct {
	Drive      resource.ID
	Collection SeqNum
}

// Error returns a string representation of the error.
func (e StatesTruncated) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d contains an inconsistent view of its states", e.Drive, e.Collection)
}

// PageNotFound reports that a requested page does not exist
// within the collection.
type PageNotFound struct {
	Drive      resource.ID
	Collection SeqNum
	Page       page.SeqNum
}

// Error returns a string representation of the error.
func (e PageNotFound) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d does not contain page %d", e.Drive, e.Collection, e.Page)
}

// PageOutOfOrder reports that a collection page could not be created
// because its page number is not the next one in the series.
type PageOutOfOrder struct {
	Drive      resource.ID
	Collection SeqNum
	Page       page.SeqNum
	Expected   page.SeqNum
}

// Error returns a string representation of the error.
func (e PageOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d page %d could not be created because it would be out of order (expected %d)", e.Drive, e.Collection, e.Page, e.Expected)
}

// PageDataInvalid reports that a requested collection contains invalid
// or unparsable page data.
type PageDataInvalid struct {
	Drive      resource.ID
	Collection SeqNum
	Page       page.SeqNum
}

// Error returns a string representation of the error.
func (e PageDataInvalid) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d contains invalid data for state %d", e.Drive, e.Collection, e.Page)
}

// PagesTruncated reports that a request for collection pages returned a
// shorter list than expected. This is an indication that the repository is
// providing an inconsistent view of its data.
type PagesTruncated struct {
	Drive      resource.ID
	Collection SeqNum
}

// Error returns a string representation of the error.
func (e PagesTruncated) Error() string {
	return fmt.Sprintf("drivestream: drive %s: collection %d contains an inconsistent view of its pages", e.Drive, e.Collection)
}
