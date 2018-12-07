package resource

import "time"

// FileData holds the properties of a file.
type FileData struct {
	Name         string    `json:"name"`
	MimeType     string    `json:"mimeType,omitempty"`
	Description  string    `json:"description,omitempty"`
	OriginalName string    `json:"originalFilename,omitempty"`
	RevisionID   string    `json:"headRevisionId"`
	MD5Checksum  string    `json:"md5Checksum"`
	Size         int64     `json:"size"`
	Created      time.Time `json:"createdTime,omitempty"`
	Modified     time.Time `json:"modifiedTime,omitempty"`
	Parents      []string  `json:"parents,omitempty"`
}

// IsDir returns true if the file data describes a directory.
func (f *FileData) IsDir() bool {
	return f.MimeType == "application/vnd.google-apps.folder"
}
