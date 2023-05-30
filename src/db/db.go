package db

import (
	"database/sql"

	"github.com/igorsfreitas/devbook-api/src/config"
	_ "github.com/lib/pq"
)

// Connect opens a connection with the database and returns it
func Connect() (*sql.DB, error) {
	db, erro := sql.Open("postgres", config.ConnectionString)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}
