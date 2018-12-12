package collection

import (
	"fmt"
	"time"

	"github.com/scjalliance/drivestream/page"
)

/*
// State is the state of a collection at a point in time.
type State struct {
	Time     time.Time `json:"time"`
	Instance string    `json:"instance"`
	StateData
}

// String returns a string representation of s.
func (s State) String() string {
	return fmt.Sprintf("%d.%d", s.Phase, s.Page)
}
*/

// State describes the state of a collection at a point in time.
type State struct {
	Time     time.Time   `json:"time"`
	Instance string      `json:"instance"`
	Phase    Phase       `json:"phase"`
	Page     page.SeqNum `json:"page,omitempty"`
}

// String returns a string representation of s.
func (s State) String() string {
	return fmt.Sprintf("%d.%d", s.Phase, s.Page)
}
