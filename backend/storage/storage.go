package storage

import (
	"zebra-rss/entries"
	"zebra-rss/sources"
	"zebra-rss/users"
)

type Storage struct {
	Users   users.Repository
	Sources sources.Repository
	Entries entries.Repository
}
