package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scjalliance/drivestream"
	"github.com/scjalliance/drivestream/collection"
	"github.com/scjalliance/drivestream/commit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func dump(ctx context.Context, app *kingpin.Application, repo drivestream.Repository, kinds []string, wanted []string) {
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

		for _, kind := range kinds {
			if ctx.Err() != nil {
				return
			}
			switch kind {
			case "collections", "collection", "cols", "col":
				cursor, err := collection.NewCursor(drv.Collections())
				if err != nil {
					app.Fatalf("failed to create collection cursor for repository %s: %v", driveID, err)
				}
				for cursor.First(); cursor.Valid(); cursor.Next() {
					reader, err := cursor.Reader()
					if err != nil {
						app.Fatalf("failed to create collection reader for repository %s: %v", driveID, err)
					}

					data, err := reader.Data()
					if err != nil {
						app.Fatalf("failed to read collection data from repository %s: %v", driveID, err)
					}

					{
						b, err := json.Marshal(data)
						if err != nil {
							fmt.Printf("%s: COLLECTION %d: DATA: parse error: %v\n", prefix, cursor.SeqNum(), err)
						} else {
							fmt.Printf("%s: COLLECTION %d: DATA: %s\n", prefix, cursor.SeqNum(), string(b))
						}
					}

					states, err := reader.States()
					if err != nil {
						app.Fatalf("failed to read collection states from repository %s: %v", driveID, err)
					}
					for i, state := range states {
						b, err := json.Marshal(state)
						if err != nil {
							fmt.Printf("%s: COLLECTION %d: STATE %d: parse error: %v\n", prefix, cursor.SeqNum(), i, err)
						} else {
							fmt.Printf("%s: COLLECTION %d: STATE %d: %v\n", prefix, cursor.SeqNum(), i, string(b))
						}
					}

					pages, err := reader.Pages()
					if err != nil {
						app.Fatalf("failed to read pages from repository %s: %v", driveID, err)
					}
					for i, pg := range pages {
						if ctx.Err() != nil {
							return
						}
						for c, change := range pg.Changes {
							b, err := json.Marshal(change)
							if err != nil {
								fmt.Printf("%s: COLLECTION %d: PAGE %d: CHANGE %d: parse error: %v\n", prefix, cursor.SeqNum(), i, c, err)
							} else {
								fmt.Printf("%s: COLLECTION %d: PAGE %d: CHANGE %d: %v\n", prefix, cursor.SeqNum(), i, c, string(b))
							}
						}
					}
				}

			case "commits", "commit", "com":
				cursor, err := commit.NewCursor(drv.Commits())
				if err != nil {
					app.Fatalf("failed to create commit cursor for repository %s: %v", driveID, err)
				}
				for cursor.First(); cursor.Valid(); cursor.Next() {
					reader, err := cursor.Reader()
					if err != nil {
						app.Fatalf("failed to create commit reader for repository %s: %v", driveID, err)
					}

					data, err := reader.Data()
					if err != nil {
						app.Fatalf("failed to read commit data from repository %s: %v", driveID, err)
					}

					{
						b, err := json.Marshal(data)
						if err != nil {
							fmt.Printf("%s: COMMIT %d: DATA: parse error: %v\n", prefix, cursor.SeqNum(), err)
						} else {
							fmt.Printf("%s: COMMIT %d: DATA: %s\n", prefix, cursor.SeqNum(), string(b))
						}
					}

					states, err := reader.States()
					if err != nil {
						app.Fatalf("failed to read commit states from repository %s: %v", driveID, err)
					}
					for i, state := range states {
						b, err := json.Marshal(state)
						if err != nil {
							fmt.Printf("%s: COMMIT %d: STATE %d: parse error: %v\n", prefix, cursor.SeqNum(), i, err)
						} else {
							fmt.Printf("%s: COMMIT %d: STATE %d: %v\n", prefix, cursor.SeqNum(), i, string(b))
						}
					}
				}
			}
		}

		//fmt.Printf("%s: %s\n", prefix, teamDrive.Name)
	}

	printMemUsage()
}
