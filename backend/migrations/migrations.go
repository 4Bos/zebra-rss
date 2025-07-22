package migrations

import (
	"database/sql"
)

var migrations = []migration{
	{version: 00001, controller: createSourcesTable},
	{version: 00002, controller: createEntriesTable},
}

type migration struct {
	version    int
	controller controller
}

type controller struct {
	Up   func(db *sql.DB) error
	Down func(db *sql.DB) error
}

func Migrate(db *sql.DB) error {
	if err := createOptionsTable(db); err != nil {
		return err
	}

	version, err := getCurrentVersion(db)

	if err != nil {
		return err
	}

	for _, m := range migrations {
		if version >= m.version {
			continue
		}

		err = m.controller.Up(db)

		if err == nil {
			err = setCurrentVersion(db, m.version)
		}

		if err != nil {
			break
		}
	}

	return err
}

func createOptionsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS zebra.options (
		key VARCHAR(64) NOT NULL,
		value VARCHAR,

		CONSTRAINT options_pkey PRIMARY KEY (key)
	)`

	_, err := db.Exec(query)

	return err
}

func getCurrentVersion(db *sql.DB) (int, error) {
	query := `SELECT value FROM zebra.options WHERE key = 'version'`

	var version int

	err := db.QueryRow(query).Scan(&version)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return version, nil
}

func setCurrentVersion(db *sql.DB, version int) error {
	query := `SELECT value FROM zebra.options WHERE key = 'version'`

	var currentVersion int

	err := db.QueryRow(query).Scan(&currentVersion)

	switch err {
	case nil:
		_, err = db.Exec("UPDATE zebra.options SET value = $1 WHERE key = 'version'", version)
	case sql.ErrNoRows:
		_, err = db.Exec("INSERT INTO zebra.options (key, value) VALUES ('version', $1)", version)
	}

	return err
}
