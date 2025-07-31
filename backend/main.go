package main

import (
	"fmt"
	"log"
	"time"
	"zebra-rss/db"
	"zebra-rss/entries"
	"zebra-rss/migrations"
	"zebra-rss/sources"
	"zebra-rss/storage"
	"zebra-rss/updater"
	"zebra-rss/users"
	"zebra-rss/web"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Connect to the database.
	opts, err := db.CollectOptions()

	if err != nil {
		log.Fatal(err)
	}

	conn := db.Connet(opts)

	// Run migrations.
	if err := migrations.Migrate(conn); err != nil {
		log.Fatal(err)
	}

	// Initialize storage.
	storage := &storage.Storage{
		Users:   users.NewRepository(conn),
		Sources: sources.NewRepository(conn),
		Entries: entries.NewRepository(conn),
	}

	// Run webserver.
	go web.StartServer(storage)

	// Main cycle.
	for {
		sourcesToUpdate, err := storage.Sources.GetSourcesToUpdate()

		if err != nil {
			log.Fatal(err)
		}

		for _, source := range sourcesToUpdate {
			err := updater.UpdateSource(conn, storage.Sources, storage.Entries, source)

			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
