package memrepo

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
)

// Collection holds data for a collection.
type Collection struct {
	Data   collection.Data
	States []collection.State
	Pages  []page.Data
}
