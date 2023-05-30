package repositories

import (
	"database/sql"

	"github.com/igorsfreitas/devbook-api/src/models"
)

// Users represents a user repository
type Users struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *Users {
	return &Users{db}
}

// Create creates a new user
func (repository Users) Create(user models.User) (uint64, error) {
	statement, err := repository.db.Prepare("insert into users (name, nick, email, password) values ($1, $2, $3, $4) returning id")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	err = statement.QueryRow(user.Name, user.Nick, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return 0, err
	}

	return uint64(user.ID), nil

}
