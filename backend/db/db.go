package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type options struct {
	host string
	port string
	user string
	pass string
	name string
}

// Establishes a connection to the database.
func Connet(opts options) *sql.DB {
	origin := opts.host

	if len(opts.port) != 0 {
		origin += ":" + opts.port
	}

	connStr := "postgres://" + opts.user + ":" + opts.pass + "@" + origin + "/" + opts.name + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Collects option values from environment variables.
func CollectOptions() (options, error) {
	var result options

	optDefs := []struct {
		env      string
		option   *string
		required bool
	}{
		{env: "DB_HOST", option: &result.host, required: true},
		{env: "DB_PORT", option: &result.port, required: false},
		{env: "DB_USER", option: &result.user, required: true},
		{env: "DB_PASS", option: &result.pass, required: true},
		{env: "DB_NAME", option: &result.name, required: true},
	}

	for i, optDef := range optDefs {
		value, exists := os.LookupEnv(optDef.env)

		if optDef.required && (!exists || len(value) == 0) {
			return result, errors.New("environment variable " + optDef.env + " is required")
		}

		*optDefs[i].option = value
	}

	return result, nil
}
