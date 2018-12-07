package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/driveapicollector"
	drive "google.golang.org/api/drive/v3"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func update(ctx context.Context, app *kingpin.Application, db *DB, includeStats bool, email string, interval time.Duration, wanted []string) {
	if ctx.Err() != nil {
		return
	}

	client := getClient(getConfig(drive.DriveReadonlyScope))
	driveService, err := drive.New(client)
	if err != nil {
		app.Fatalf("failed to create google drive client: %v", err)
	}

	for {
		if ctx.Err() != nil {
			return
		}

		selection, err := selectTeamDrives(ctx, driveService, email, wanted)
		if err != nil {
			app.Fatalf("failed to enumerate team drives: %v", err)
		}

		for _, teamDrive := range selection {
			if ctx.Err() != nil {
				return
			}

			prefix := fmt.Sprintf("DRIVE %s", teamDrive.ID)

			fmt.Printf("%s: NAME: %s\n", prefix, teamDrive.Name)

			repo, existing := db.Repository(teamDrive.ID)
			if !existing {
				fmt.Printf("%s: INIT: Repository (%s)\n", prefix, db.Kind())
			}

			collector := driveapicollector.New(driveService, string(teamDrive.ID))
			stream := drivestream.New(repo, drivestream.WithLogger(os.Stdout))
			stream.Update(ctx, collector)
		}

		if includeStats {
			printMemUsage()
			// TODO: Include database statistics
		}

		if interval == 0 {
			return
		}

		if ctx.Err() != nil {
			return
		}

		fmt.Printf("Sleeping %s\n", interval)

		t := time.NewTimer(interval)
		select {
		case <-t.C:
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}
			return
		}
	}
}
