package commit

import (
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/page"
)

// Source identifies the source of a commit within a collection.
type Source struct {
	Collection collection.SeqNum `json:"collection"`
	Page       page.SeqNum       `json:"page,omitempty"`
	Index      int               `json:"index,omitempty"`
}
