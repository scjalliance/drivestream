package fileview

import (
	"github.com/scjalliance/drivestream/commit"
	"github.com/scjalliance/drivestream/resource"
)

// Data holds data about a file view.
type Data struct {
	File    resource.ID
	Drive   resource.ID
	Commit  commit.SeqNum
	Version resource.Version
}
