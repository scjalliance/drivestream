package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/scjalliance/drivestream/driveapicollector"
	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

func selectTeamDrives(ctx context.Context, s *drive.Service, email string, wanted []string) (drives []resource.Drive, err error) {
	var token string
	for {
		call := s.Teamdrives.List()
		call.Context(ctx)
		call.Fields("nextPageToken", "teamDrives(id,name,capabilities,createdTime)")
		if token != "" {
			call.PageToken(token)
		}

		list, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve team drive list: %v", err)
		}

		for _, teamDrive := range list.TeamDrives {
			if !isWanted(wanted, teamDrive.Id, teamDrive.Name) {
				continue
			}

			switch {
			case teamDrive.Capabilities == nil:
			case !teamDrive.Capabilities.CanListChildren:
			//case !teamDrive.Capabilities.CanReadRevisions:
			default:
				if perms, err := listPermissions(ctx, s, teamDrive.Id); err == nil {
					if hasDirectMembership(email, perms) {
						if record, err := driveapicollector.MarshalDrive(teamDrive); err == nil {
							drives = append(drives, record)
						}
					}
				}
			}
		}

		if list.NextPageToken == "" {
			return drives, nil
		}

		token = list.NextPageToken
	}
}

func isWanted(wanted []string, values ...string) bool {
	if len(wanted) == 0 {
		return true
	}

	for i := range wanted {
		for j := range values {
			if strings.EqualFold(wanted[i], values[j]) {
				return true
			}
		}
	}

	return false
}

func hasDirectMembership(email string, perms []resource.Permission) bool {
	for _, perm := range perms {
		if perm.Deleted {
			continue
		}
		if perm.EmailAddress != email {
			continue
		}
		return true
	}
	return false
}
