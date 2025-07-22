package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		feedPath string
		expected Rss
	}{
		{
			name:     "Basic feed",
			feedPath: "testdata/default.rss",
			expected: Rss{
				Channel: Channel{
					Title: "Channel Title",
					Items: []Item{
						{
							Title:       "Entry One",
							Link:        "https://example.com/one",
							Description: "Description One",
						},
						{
							Title:       "Entry Two",
							Link:        "https://example.com/two",
							Description: "Description Two",
						},
						{
							Title:       "Entry Three",
							Link:        "https://example.com/three",
							Description: "Description Three",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, _ := os.Open("testdata/default.rss")

			defer f.Close()

			actual, _ := Parse(f)

			assert.Equal(t, tc.expected.Channel.Title, actual.Channel.Title, "Channel title should match expected value")
			assert.Equal(t, len(tc.expected.Channel.Items), len(actual.Channel.Items), "Number of items should match expected count")

			for i, expectedItem := range tc.expected.Channel.Items {
				assert.Equal(t, expectedItem.Title, actual.Channel.Items[i].Title, "Item title at index %d should match", i)
				assert.Equal(t, expectedItem.Link, actual.Channel.Items[i].Link, "Item link at index %d should match", i)
				assert.Equal(t, expectedItem.Description, actual.Channel.Items[i].Description, "Item description at index %d should match", i)
			}
		})
	}
}

func TestParsePubDate_Success(t *testing.T) {
	testCases := []struct {
		name     string
		pubDate  string
		expected string
	}{
		{
			name:     "With timezone abbreviation",
			pubDate:  "Sat, 30 Oct 2010 23:08:27 MSK",
			expected: "2010-10-30 20:08:27 +0000 UTC",
		},
		{
			name:     "With timezone offset",
			pubDate:  "Sat, 30 Oct 2010 23:08:27 +0300",
			expected: "2010-10-30 20:08:27 +0000 UTC",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParsePubDate(tc.pubDate)

			assert.Nil(t, err, "Should not return error for valid date format")
			assert.Equal(t, tc.expected, result.String(), "Parsed date should match expected value")
		})
	}
}

func TestParsePubDate_Error(t *testing.T) {
	testCases := []struct {
		name    string
		pubDate string
	}{
		{
			name:    "Completely invalid fromat",
			pubDate: "Invalid Date",
		},
		{
			name:    "Empty string",
			pubDate: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePubDate(tc.pubDate)

			assert.NotNil(t, err, "Should return error for invalid date format")
		})
	}
}
