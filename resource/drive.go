package resource

// Drive holds information about a team drive.
type Drive struct {
	ID      ID      `json:"id"`
	Version Version `json:"version"`
	DriveData
}
