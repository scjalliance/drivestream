package resource

import (
	"time"
)

// Change holds information about a change to a team drive file.
type Change struct {
	Type    Type      `json:"type"`
	Time    time.Time `json:"time"`
	Removed bool      `json:"removed,omitempty"`
	File    `json:"file,omitempty"`
	Drive   `json:"drive,omitempty"`
}
