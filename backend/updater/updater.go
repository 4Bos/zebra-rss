package updater

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
	"zebra-rss/entries"
	"zebra-rss/parser"
	"zebra-rss/sources"
)

func UpdateSource(
	db *sql.DB,
	sourcesRepo sources.Repository,
	entriesRepo entries.Repository,
	source sources.Source,
) error {
	fmt.Printf("The source updating [%d]: %s", source.Id, source.Url)

	resp, err := http.Get(source.Url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("expected status code 200")
	}

	rss, err := parser.Parse(resp.Body)

	if err != nil {
		return err
	}

	// Update the source title.
	if err := sourcesRepo.SetTitle(source.Id, rss.Channel.Title); err != nil {
		return err
	}

	for _, item := range rss.Channel.Items {
		exists, err := entriesRepo.ExistsByHash(item.Hash)

		if err != nil {
			fmt.Println(err)
		} else if !exists {
			entry := entries.Entry{
				SourceId:    source.Id,
				Hash:        item.Hash,
				Title:       item.Title,
				Url:         item.Link,
				Content:     &item.Description,
				PublishedAt: (time.Time)(item.PubDate),
			}

			_, err := entriesRepo.Create(&entry)

			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// Mark source as scanned.
	return sourcesRepo.MarkSourceAsScanned(source.Id)
}
