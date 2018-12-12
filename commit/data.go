package commit

import "time"

// Data holds data about a commit.
type Data struct {
	Source Source
	Time   time.Time
}
