package driveapicollector

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

// MarshalDrive marshals the given team drive as a resource.
func MarshalDrive(d *drive.TeamDrive) (resource.Drive, error) {
	created, err := parseRFC3339(d.CreatedTime)
	if err != nil {
		return resource.Drive{}, fmt.Errorf("invalid creation time: %v", err)
	}

	return resource.Drive{
		ID: resource.ID(d.Id),
		DriveData: resource.DriveData{
			Name:    d.Name,
			Created: created,
		},
	}, nil
}
