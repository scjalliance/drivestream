package driveapicollector

import "time"

func parseRFC3339(v string) (t time.Time, err error) {
	if v == "" {
		return
	}
	return time.Parse(time.RFC3339, v)
}
