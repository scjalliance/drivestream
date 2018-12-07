package driveapicollector

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

// MarshalFile marshals the given file as a resource.
func MarshalFile(file *drive.File) (resource.File, error) {
	created, err := parseRFC3339(file.CreatedTime)
	if err != nil {
		return resource.File{}, fmt.Errorf("invalid creation time: %v", err)
	}

	modified, err := parseRFC3339(file.ModifiedTime)
	if err != nil {
		return resource.File{}, fmt.Errorf("invalid modification time: %v", err)
	}

	return resource.File{
		ID:      resource.ID(file.Id),
		Version: resource.Version(file.Version),
		FileData: resource.FileData{
			Name:         file.Name,
			MimeType:     file.MimeType,
			Description:  file.Description,
			OriginalName: file.OriginalFilename,
			RevisionID:   file.HeadRevisionId,
			MD5Checksum:  file.Md5Checksum,
			Size:         file.Size,
			Created:      created,
			Modified:     modified,
			Parents:      file.Parents,
		},
	}, nil
}
