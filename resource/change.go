package resource

import (
	"encoding/json"
	"fmt"
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

// MarshalJSON marshals the change as JSON-encoded bytes.
func (c Change) MarshalJSON() ([]byte, error) {
	switch c.Type {
	case TypeFile:
		return json.Marshal(struct {
			Type    Type      `json:"type"`
			Time    time.Time `json:"time"`
			Removed bool      `json:"removed,omitempty"`
			File    `json:"file"`
		}{
			Type:    c.Type,
			Time:    c.Time,
			Removed: c.Removed,
			File:    c.File,
		})
	case TypeDrive:
		return json.Marshal(struct {
			Type    Type      `json:"type"`
			Time    time.Time `json:"time"`
			Removed bool      `json:"removed,omitempty"`
			Drive   `json:"drive"`
		}{
			Type:    c.Type,
			Time:    c.Time,
			Removed: c.Removed,
			Drive:   c.Drive,
		})
	default:
		return nil, fmt.Errorf("unknown resource type \"%s\"", c.Type)
	}
}
