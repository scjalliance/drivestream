package commit

import (
	"fmt"
)

// NotFound reports that a requested commit does not exist
// within the repository.
type NotFound struct {
	SeqNum SeqNum
}

// Error returns a string representation of the error.
func (e NotFound) Error() string {
	return fmt.Sprintf("drivestream commit %d could not be found", e.SeqNum)
}

// OutOfOrder reports that a commit could not be created because its
// sequence number is not the next one in the series.
type OutOfOrder struct {
	SeqNum   SeqNum
	Expected SeqNum
}

// Error returns a string representation of the error.
func (e OutOfOrder) Error() string {
	return fmt.Sprintf("drivestream commit %d could not be created because it would be out of order (expected %d)", e.SeqNum, e.Expected)
}

// InvalidData reports that a requested commit contains invalid
// or unparsable data.
type InvalidData struct {
	SeqNum SeqNum
}

// Error returns a string representation of the error.
func (e InvalidData) Error() string {
	return fmt.Sprintf("drivestream commit %d contains invalid data", e.SeqNum)
}

// StateNotFound reports that a requested state does not exist
// within the commit.
type StateNotFound struct {
	SeqNum   SeqNum
	StateNum StateNum
}

// Error returns a string representation of the error.
func (e StateNotFound) Error() string {
	return fmt.Sprintf("drivestream commit %d does not contain state %d", e.SeqNum, e.StateNum)
}

// StateOutOfOrder reports that a commit state could not be created
// because its state number is not the next one in the series.
type StateOutOfOrder struct {
	SeqNum   SeqNum
	StateNum StateNum
	Expected StateNum
}

// Error returns a string representation of the error.
func (e StateOutOfOrder) Error() string {
	return fmt.Sprintf("drivestream commit %d state %d could not be created because it would be out of order (expected %d)", e.SeqNum, e.StateNum, e.Expected)
}

// InvalidState reports that a requested commit contains invalid
// or unparsable state.
type InvalidState struct {
	SeqNum   SeqNum
	StateNum StateNum
}

// Error returns a string representation of the error.
func (e InvalidState) Error() string {
	return fmt.Sprintf("drivestream commit %d contains invalid data for state %d", e.SeqNum, e.StateNum)
}

// TruncatedStates reports that a request for commit states returned a
// shorter list than expected. This is an indication that the repository is
// providing an inconsistent view of its data.
type TruncatedStates struct {
	SeqNum SeqNum
}

// Error returns a string representation of the error.
func (e TruncatedStates) Error() string {
	return fmt.Sprintf("drivestream commit %d contains an inconsistent view of its states", e.SeqNum)
}
