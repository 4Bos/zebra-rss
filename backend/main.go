package main

import (
	"fmt"
	"log"
	"time"
	"zebra-rss/db"
	"zebra-rss/entries"
	"zebra-rss/migrations"
	"zebra-rss/sources"
	"zebra-rss/updater"

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

	// Initialize repos.
	sourcesRepo := sources.NewRepository(conn)
	entriesRepo := entries.NewRepository(conn)

	// Main cycle.
	for {
		sourcesToUpdate, err := sourcesRepo.GetSourcesToUpdate()

		if err != nil {
			log.Fatal(err)
		}

		for _, source := range sourcesToUpdate {
			err := updater.UpdateSource(conn, sourcesRepo, entriesRepo, source)

			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
