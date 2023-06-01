package repositories

import (
	"database/sql"
	"fmt"

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

// Find finds all users filtered by nick or name
func (repository Users) Find(nickOrName string) ([]models.User, error) {
	nickOrName = fmt.Sprintf("%%%s%%", nickOrName) // %nickOrName%

	linhas, err := repository.db.Query(
		"select id, name, nick, email, created_at from users where lower(name) like $1 or lower(nick) like $2",
		nickOrName,
		nickOrName,
	)

	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var users []models.User
	for linhas.Next() {
		var user models.User

		if err = linhas.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

// FindByID finds a user by id
func (repository Users) FindByID(userID uint64) (models.User, error) {
	linha, err := repository.db.Query(
		"select id, name, nick, email, created_at from users where id = $1",
		userID,
	)

	if err != nil {
		return models.User{}, err
	}
	defer linha.Close()

	var user models.User
	if linha.Next() {
		if err = linha.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// Update updates a user
func (repository Users) Update(userID uint64, user models.User) error {
	statement, err := repository.db.Prepare(
		"update users set name = $1, nick = $2, email = $3 where id = $4",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nick, user.Email, userID); err != nil {
		return err
	}

	return nil

}

// Delete deletes a user
func (repository Users) Delete(userID uint64) error {
	statement, err := repository.db.Prepare("delete from users where id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID); err != nil {
		return err
	}

	return nil

}

// FindByEmail finds a user by email and returns it with id and password with hash
func (repository Users) FindByEmail(email string) (models.User, error) {
	line, err := repository.db.Query("select id, password from users where email = $1", email)
	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User
	if line.Next() {
		if err = line.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

// Follow follows a user
func (repository Users) FollowUser(userID, followerID uint64) error {
	statement, err := repository.db.Prepare("insert into followers (user_id, follower_id) values ($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

// Unfollow unfollows a user
func (repository Users) UnfollowUser(userID, followerID uint64) error {
	statement, err := repository.db.Prepare("delete from followers where user_id = $1 and follower_id = $2")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

// FindFollowers finds all followers of a user
func (repository Users) FindFollowers(userID uint64) ([]models.User, error) {
	lines, err := repository.db.Query(`
		select u.id, u.name, u.nick, u.email, u.created_at
			from users u 
			inner join followers f on u.id = f.follower_id
			where f.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// FindFollowing finds all users that a user is following
func (repository Users) FindFollowing(userID uint64) ([]models.User, error) {
	lines, err := repository.db.Query(`
		select u.id, u.name, u.nick, u.email, u.created_at
			from users u 
			inner join followers f on u.id = f.user_id
			where f.follower_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var users []models.User
	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Update password updates a user password
func (repository Users) UpdatePassword(userID uint64, password string) error {
	statement, err := repository.db.Prepare("update users set password = $1 where id = $2")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userID); err != nil {
		return err
	}

	return nil
}

// GetUserPassword gets a user password
func (repository Users) GetUserPassword(userID uint64) (string, error) {
	line, err := repository.db.Query("select password from users where id = $1", userID)
	if err != nil {
		return "", err
	}
	defer line.Close()

	var user models.User
	if line.Next() {
		if err = line.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}
