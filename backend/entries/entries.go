package entries

import (
	"database/sql"
	"time"
)

type Entry struct {
	Id          int64
	SourceId    int64
	Hash        string
	Title       string
	Url         string
	Content     *string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	ExistsByHash(hash string) (bool, error)
	Create(entry *Entry) (*Entry, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ExistsByHash(hash string) (bool, error) {
	rows, err := r.db.Query("SELECT id FROM zebra.entries WHERE hash = $1", hash)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (r *repository) Create(entry *Entry) (*Entry, error) {
	query := "INSERT INTO zebra.entries (source_id, hash, title, url, content, published_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	row := r.db.QueryRow(query, entry.SourceId, entry.Hash, entry.Title, entry.Url, entry.Content, entry.PublishedAt)
	err := row.Scan(&entry.Id)

	if err != nil {
		return nil, err
	}

	return entry, nil
}
