package page

import (
	"time"

	"github.com/scjalliance/drivestream/resource"
)

// Data holds a page of collected changes.
type Data struct {
	Type           Type
	Collected      time.Time
	PageToken      string
	NextPageToken  string
	NextStartToken string
	Changes        []resource.Change
}

// Last returns true if the page is the last one of its type.
func (d *Data) Last() bool {
	return d.NextPageToken == ""
}
