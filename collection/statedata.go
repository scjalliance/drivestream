package collection

import (
	"fmt"

	"github.com/scjalliance/drivestream/page"
)

// StateData holds data about a collection's state.
type StateData struct {
	Phase Phase       `json:"phase"`
	Page  page.SeqNum `json:"page,omitempty"`
}

// String returns a string representation of s.
func (s StateData) String() string {
	return fmt.Sprintf("%d.%d", s.Phase, s.Page)
}
