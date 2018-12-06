package collection

import (
	"fmt"
	"time"
)

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
