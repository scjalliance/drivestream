package main

import (
	"context"
	"fmt"

	"github.com/scjalliance/drivestream/driveapicollector"
	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

func listPermissions(ctx context.Context, s *drive.Service, id string) (perms []resource.Permission, err error) {
	var token string

	for {
		call := s.Permissions.List(id)
		call.Context(ctx)
		call.SupportsTeamDrives(true)
		call.UseDomainAdminAccess(true)
		call.Fields("nextPageToken", "permissions(id,type,emailAddress,domain,role,displayName,expirationTime,deleted)")

		if token != "" {
			call.PageToken(token)
		}

		list, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve teamdrive permissions: %v", err)
		}

		for i, perm := range list.Permissions {
			record, err := driveapicollector.MarshalPermission(perm)
			if err != nil {
				return nil, fmt.Errorf("permission list parsing failed: record %d: %v", i, err)
			}
			perms = append(perms, record)
		}

		if list.NextPageToken == "" {
			return perms, nil
		}

		token = list.NextPageToken
	}
}
