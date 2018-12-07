package commit

import "time"

// State is the state of a collection at a point in time.
type State struct {
	Time     time.Time `json:"time"`
	Instance string    `json:"instance"`
	StateData
}
