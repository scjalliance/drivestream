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
		switch {
		case change.Removed || change.File == nil:
			return resource.Change{
				Type:    resource.TypeFile,
				Time:    changed,
				Removed: true,
				File: resource.File{
					ID: resource.ID(change.File.Id),
				},
			}, nil
		default:
			record, err := MarshalFile(change.File)
			if err != nil {
				return resource.Change{}, err
			}
			return resource.Change{
				Type: resource.TypeFile,
				Time: changed,
				File: record,
			}, nil
		}
	case "teamDrive":
		switch {
		case change.Removed || change.TeamDrive == nil:
			return resource.Change{
				Type:    resource.TypeDrive,
				Time:    changed,
				Removed: true,
				Drive: resource.Drive{
					ID: resource.ID(change.TeamDriveId),
				},
			}, nil
		default:
			record, err := MarshalDrive(change.TeamDrive)
			if err != nil {
				return resource.Change{}, err
			}
			return resource.Change{
				Type:  resource.TypeDrive,
				Time:  changed,
				Drive: record,
			}, nil
		}
	default:
		return resource.Change{}, fmt.Errorf("unknown change type: \"%s\"", change.Type)
	}
}
