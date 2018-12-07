package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scjalliance/drivestream/collection"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func dump(ctx context.Context, app *kingpin.Application, db *DB, kinds []string, wanted []string) {
	if ctx.Err() != nil {
		return
	}

	ids, err := db.Enumerate()
	if err != nil {
		app.Fatalf("failed to enumerate drivestream database: %v", err)
	}

	for _, teamDriveID := range ids {
		repo, _ := db.Repository(teamDriveID)

		prefix := fmt.Sprintf("DRIVE %s", teamDriveID)

		if data, ok := driveData(repo); ok {
			if !isWanted(wanted, string(teamDriveID), data.Name) {
				continue
			}
			b, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("%s: DATA: parse error: %v\n", prefix, err)
			} else {
				fmt.Printf("%s: DATA: %v\n", prefix, string(b))
			}
		} else {
			if !isWanted(wanted, string(teamDriveID)) {
				continue
			}
		}

		for _, kind := range kinds {
			if ctx.Err() != nil {
				return
			}
			switch kind {
			case "collections", "collection", "cols", "col":
				cursor, err := collection.NewCursor(repo)
				if err != nil {
					app.Fatalf("failed to create collection cursor for repository %s: %v", teamDriveID, err)
				}
				for cursor.First(); cursor.Valid(); cursor.Next() {
					reader, err := cursor.Reader()
					if err != nil {
						app.Fatalf("failed to create collection reader for repository %s: %v", teamDriveID, err)
					}

					data, err := reader.Data()
					if err != nil {
						app.Fatalf("failed to read data from repository %s: %v", teamDriveID, err)
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
						app.Fatalf("failed to read states from repository %s: %v", teamDriveID, err)
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
						app.Fatalf("failed to read pages from repository %s: %v", teamDriveID, err)
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
			}
		}

		//fmt.Printf("%s: %s\n", prefix, teamDrive.Name)
	}

	printMemUsage()
}
