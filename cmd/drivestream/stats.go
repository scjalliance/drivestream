package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scjalliance/drivestream"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func stats(ctx context.Context, app *kingpin.Application, repo drivestream.Repository, selections []string, wanted []string) {
	if ctx.Err() != nil {
		return
	}

	ids, err := repo.Drives().List()
	if err != nil {
		app.Fatalf("failed to enumerate drivestream database: %v", err)
	}

	for _, driveID := range ids {
		drv := repo.Drive(driveID)

		prefix := fmt.Sprintf("DRIVE %s", driveID)

		if data, ok := driveData(drv); ok {
			if !isWanted(wanted, string(driveID), data.Name) {
				continue
			}
			b, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("%s: DATA: parse error: %v\n", prefix, err)
			} else {
				fmt.Printf("%s: DATA: %v\n", prefix, string(b))
			}
		} else {
			if !isWanted(wanted, string(driveID)) {
				continue
			}
		}
		stats, err := drv.Stats()
		if err != nil {
			app.Fatalf("failed to collect statistics for drive %s: %v", driveID, err)
		}
		const line = "--------"
		fmt.Printf("%s\n%s\n%s\n", line, strings.Join(stats.Summary(), "\n"), line)
	}
}
