package resource

import "time"

// Permission holds a team drive permission entry.
type Permission struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	EmailAddress string    `json:"email"`
	Domain       string    `json:"domain,omitempty"`
	Role         string    `json:"role,omitempty"`
	DisplayName  string    `json:"displayName,omitempty"`
	Expiration   time.Time `json:"expiration,omitempty"`
	Deleted      bool      `json:"deleted,omitempty"`
}
