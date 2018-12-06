package resource

// UserData holds information about a team drive user at a point in time.
type UserData struct {
	DisplayName  string `json:"displayName"`
	PermissionID string `json:"permission"`
	EmailAddress string `json:"email"`
}
