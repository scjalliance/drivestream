package resource

// File holds information about a file.
type File struct {
	ID      ID      `json:"id"`
	Version Version `json:"version"`
	FileData
}
