package sources

import (
	"database/sql"
	"time"
)

type Source struct {
	Id        int64
	Title     *string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	ScannedAt *time.Time
}

type Repository interface {
	GetSourcesToUpdate() ([]Source, error)
	MarkSourceAsScanned(sourceId int64) error
	SetTitle(sourceId int64, title string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetSourcesToUpdate() ([]Source, error) {
	var source Source
	var result []Source

	rows, err := r.db.Query("SELECT id, title, url, created_at, updated_at, scanned_at FROM zebra.sources WHERE scanned_at IS NULL OR scanned_at < NOW() - interval '1'")

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&source.Id, &source.Title, &source.Url, &source.CreatedAt, &source.UpdatedAt, &source.ScannedAt)

		result = append(result, source)
	}

	return result, err
}

func (r *repository) MarkSourceAsScanned(sourceId int64) error {
	_, err := r.db.Exec("UPDATE zebra.sources SET scanned_at = NOW() WHERE id = $1", sourceId)

	return err
}

func (r *repository) SetTitle(sourceId int64, title string) error {
	_, err := r.db.Exec("UPDATE zebra.sources SET title = $1 WHERE id = $2", title, sourceId)

	return err
}
