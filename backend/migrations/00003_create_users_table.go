package migrations

import "database/sql"

var createUsersTable = controller{
	Up: func(db *sql.DB) error {
		query := `CREATE TABLE zebra.users (
			id BIGSERIAL NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

			CONSTRAINT users__id__pkey PRIMARY KEY (id),
    		CONSTRAINT users__email__ukey UNIQUE (email)
		)`

		_, err := db.Exec(query)

		return err
	},
	Down: func(db *sql.DB) error {
		_, err := db.Exec("DROP TABLE zebra.users")

		return err
	},
}
