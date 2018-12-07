package driveapicollector

import (
	"context"
	"fmt"

	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

// A Collector is responsible for collecting team drive file data from a
// drive service.
//
// Collectors should be created by calling New.
type Collector struct {
	id      string
	service *drive.Service
}

// New returns a new collector for the requested team drive.
func New(s *drive.Service, teamDriveID string) *Collector {
	return &Collector{
		id:      teamDriveID,
		service: s,
	}
}

// ChangeToken returns the starting token for a new stream of changes.
func (c *Collector) ChangeToken(ctx context.Context) (startToken string, err error) {
	call := c.service.Changes.GetStartPageToken()
	call.Context(ctx)
	call.SupportsTeamDrives(true)
	call.TeamDriveId(c.id)

	result, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("failed to get starting token for change list: %v", err)
	}

	return result.StartPageToken, nil
}

// Drive collects the current drive data, formatted in the same manner as
// a change.
func (c *Collector) Drive(ctx context.Context) (resource.Change, error) {
	if err := ctx.Err(); err != nil {
		return resource.Change{}, err
	}

	call := c.service.Teamdrives.Get(c.id)
	call.Context(ctx)
	call.Fields("id,name,createdTime")

	result, err := call.Do()
	if err != nil {
		return resource.Change{}, fmt.Errorf("drive get call failed: %v", err)
	}

	record, err := MarshalDrive(result)
	if err != nil {
		return resource.Change{}, err
	}

	return resource.Change{
		Type:  resource.TypeDrive,
		Time:  record.Created,
		Drive: record,
	}, nil
}

// Files collects a set of files into p, starting from the file identified
// by token. It returns the number of files collected as n, and returns a
// non-empty nextToken if there are additional files in the list yet to be
// read.
//
// If the provided token is empty it will start at the first file within
// the team drive.
//
// If the length of p is zero Files will panic.
func (c *Collector) Files(ctx context.Context, token string, p []resource.Change) (n int, nextToken string, err error) {
	bufferSize := len(p)
	if bufferSize == 0 {
		if p == nil {
			panic("unable to collect files into nil buffer")
		}
		panic("unable to collect files into empty buffer")
	}

	for n < bufferSize {
		if err := ctx.Err(); err != nil {
			return 0, token, err
		}

		call := c.service.Files.List()
		call.Context(ctx)
		call.SupportsTeamDrives(true)
		call.IncludeTeamDriveItems(true)
		call.TeamDriveId(c.id)
		call.Corpora("teamDrive")
		call.Spaces("drive")
		call.Fields("nextPageToken", "files(id,name,mimeType,description,parents,version,createdTime,modifiedTime,lastModifyingUser,originalFilename,md5Checksum,headRevisionId,size)")
		call.PageSize(c.pageSize(bufferSize - n))
		if token != "" {
			call.PageToken(token)
		}

		result, err := call.Do()
		if err != nil {
			return n, token, fmt.Errorf("file list call failed: %v", err)
		}

		for i, file := range result.Files {
			record, err := MarshalFile(file)
			if err != nil {
				return n, token, fmt.Errorf("file list parsing failed: record %d: %v", i, err)

			}
			p[n] = resource.Change{
				Type: resource.TypeFile,
				Time: record.Modified,
				File: record,
			}
			n++
		}

		if result.NextPageToken != "" {
			token = result.NextPageToken
		} else {
			return n, "", nil
		}
	}

	return n, token, nil
}

// Changes collects a set of changes into p, up to len(p), starting from
// the change identified by token.
//
// If len(p) is zero it will panic.
//
// The number of changes collected are returned in n.
//
// If there more changes to be collected in the current set, nextToken
// will be non-empty. If there are no more changes in the current set
// then nextStartToken will hold the starting token for the next set.
func (c *Collector) Changes(ctx context.Context, token string, p []resource.Change) (n int, nextToken string, nextStartToken string, err error) {
	bufferSize := len(p)
	if bufferSize == 0 {
		if p == nil {
			panic("unable to collect changes into nil buffer")
		}
		panic("unable to collect changes into empty buffer")
	}

	nextToken = token

	for n < bufferSize {
		if err := ctx.Err(); err != nil {
			return n, nextToken, nextStartToken, err
		}

		call := c.service.Changes.List(nextToken)
		call.Context(ctx)
		call.SupportsTeamDrives(true)
		call.IncludeTeamDriveItems(true)
		call.TeamDriveId(c.id)
		call.IncludeRemoved(true)
		call.Spaces("drive")
		call.Fields("nextPageToken", "newStartPageToken", "changes(fileId,removed,time,file(id,name,mimeType,description,parents,version,createdTime,modifiedTime,lastModifyingUser,originalFilename,md5Checksum,headRevisionId,size),type,teamDriveId,teamDrive(id,name,createdTime))")
		call.PageSize(c.pageSize(bufferSize - n))

		result, err := call.Do()
		if err != nil {
			return n, nextToken, nextStartToken, fmt.Errorf("failed to retrieve change list: %v", err)
		}

		// Make sure we didn't get back a bigger page than we asked for,
		// because that would leave us with a nextToken that skips records
		if len(result.Changes) > bufferSize-n {
			return n, nextToken, nextStartToken, fmt.Errorf("drive changes API call returned a larger page than requested")
		}

		for i, change := range result.Changes {
			record, err := MarshalChange(change)
			if err != nil {
				return n, nextToken, nextStartToken, fmt.Errorf("change list parsing failed: record %d: %v", i, err)
			}
			p[n] = record
			n++
		}

		nextToken = result.NextPageToken
		nextStartToken = result.NewStartPageToken
		switch {
		case nextToken == "" && nextStartToken == "":
			return n, nextToken, nextStartToken, fmt.Errorf("failed to receive next page token")
		case nextToken == "":
			return n, nextToken, nextStartToken, nil
		}
	}

	return n, nextToken, nextStartToken, nil
}

func (c *Collector) pageSize(bufferSize int) int64 {
	const maxPageSize = 1000

	if bufferSize > maxPageSize {
		return maxPageSize
	}
	return int64(bufferSize)
}
