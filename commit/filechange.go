package commit

import "github.com/scjalliance/drivestream/resource"

// FileChange describes a file change in a commit.
type FileChange struct {
	File    resource.ID
	Version resource.Version
}
