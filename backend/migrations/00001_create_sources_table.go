package migrations

import "database/sql"

var createSourcesTable = controller{
	Up: func(db *sql.DB) error {
		query := `CREATE TABLE sources (
			id BIGSERIAL NOT NULL,
			title VARCHAR,
			url VARCHAR NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			scanned_at TIMESTAMP,

			CONSTRAINT sources_id_pkey PRIMARY KEY (id),
    		CONSTRAINT sources_url_ukey UNIQUE (url)
		)`

		_, err := db.Exec(query)

		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec("DROP TABLE sources")

		return err
	},
}
