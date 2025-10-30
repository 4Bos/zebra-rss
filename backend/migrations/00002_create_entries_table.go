package migrations

import "database/sql"

var createEntriesTable = controller{
	Up: func(db *sql.DB) error {
		query := `CREATE TABLE entries (
			id BIGSERIAL NOT NULL,
			source_id BIGINT NOT NULL,
			hash VARCHAR(32) NOT NULL,
			title VARCHAR NOT NULL,
			url VARCHAR NOT NULL,
			content VARCHAR,
			published_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

			CONSTRAINT entries_id_pkey PRIMARY KEY (id),
    		CONSTRAINT entries_source_id_hash_ukey UNIQUE (source_id, hash)
		)`

		_, err := db.Exec(query)

		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec("DROP TABLE entries")

		return err
	},
}
