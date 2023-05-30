package db

import (
	"database/sql"

	"github.com/igorsfreitas/devbook-api/src/config"
	_ "github.com/lib/pq"
)

// Connect opens a connection with the database and returns it
func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnectionString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
