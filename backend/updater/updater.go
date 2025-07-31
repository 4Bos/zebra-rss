package updater

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"zebra-rss/entries"
	"zebra-rss/parser"
	"zebra-rss/sources"
	"zebra-rss/storage"
)

func UpdateSource(
	storage *storage.Storage,
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
	if err := storage.Sources.SetTitle(source.Id, rss.Channel.Title); err != nil {
		return err
	}

	for _, item := range rss.Channel.Items {
		exists, err := storage.Entries.ExistsByHash(item.Hash)

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

			_, err := storage.Entries.Create(&entry)

			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// Mark source as scanned.
	return storage.Sources.MarkSourceAsScanned(source.Id)
}
