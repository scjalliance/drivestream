package resource

import "time"

// FileData holds the properties of a file.
type FileData struct {
	Name         string
	MimeType     string
	Description  string
	OriginalName string
	RevisionID   string
	MD5Checksum  string
	Size         int64
	Created      time.Time
	Modified     time.Time
	Parents      []string
}

// IsDir returns true if the file data describes a directory.
func (f *FileData) IsDir() bool {
	return f.MimeType == "application/vnd.google-apps.folder"
}
