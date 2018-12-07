package driveapicollector

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

// MarshalChange marshals the given change as a resource.
func MarshalChange(change *drive.Change) (resource.Change, error) {
	changed, err := parseRFC3339(change.Time)
	if err != nil {
		return resource.Change{}, fmt.Errorf("invalid change time: %v", err)
	}

	switch change.Type {
	case "file":
		if change.File == nil {
			return resource.Change{
				Type:    resource.TypeFile,
				Time:    changed,
				Removed: change.Removed,
				File: resource.File{
					ID: resource.ID(change.FileId),
				},
			}, nil
		}
		record, err := MarshalFile(change.File)
		if err != nil {
			return resource.Change{}, err
		}
		return resource.Change{
			Type:    resource.TypeFile,
			Time:    changed,
			Removed: change.Removed,
			File:    record,
		}, nil
	case "teamDrive":
		if change.TeamDrive == nil {
			return resource.Change{
				Type:    resource.TypeDrive,
				Time:    changed,
				Removed: change.Removed,
				Drive: resource.Drive{
					ID: resource.ID(change.TeamDriveId),
				},
			}, nil
		}
		record, err := MarshalDrive(change.TeamDrive)
		if err != nil {
			return resource.Change{}, err
		}
		return resource.Change{
			Type:    resource.TypeDrive,
			Time:    changed,
			Removed: change.Removed,
			Drive:   record,
		}, nil
	default:
		return resource.Change{}, fmt.Errorf("unknown change type: \"%s\"", change.Type)
	}
}
