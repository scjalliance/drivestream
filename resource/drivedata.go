package resource

import (
	"time"
)

// DriveData holds the properties of a team drive.
type DriveData struct {
	Name        string       `json:"name"`
	Created     time.Time    `json:"created"`
	Permissions []Permission `json:"permissions"`
}
