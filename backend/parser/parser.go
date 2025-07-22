package parser

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io"
	"time"
)

func Parse(r io.Reader) (*Rss, error) {
	var rss Rss

	err := xml.NewDecoder(r).Decode(&rss)

	if err != nil {
		return nil, err
	}

	// Calculate item hashes.
	for index, item := range rss.Channel.Items {
		hash := md5.New()
		hash.Write([]byte(item.Link))

		rss.Channel.Items[index].Hash = hex.EncodeToString(hash.Sum(nil))
	}

	return &rss, nil
}

// Parses the publication date and returns it in UTC.
func ParsePubDate(pubDate string) (*time.Time, error) {
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 MST",
		// Some feeds may contain dates with timezone offset instead of timezone shortcut.
		// https://www.sublimetext.com/blog/feed
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	var result time.Time
	var err error

	for _, format := range formats {
		result, err = time.Parse(format, pubDate)

		if err == nil {
			result = result.UTC()
			return &result, nil
		}
	}

	return nil, err
}
