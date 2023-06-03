package repositories

import (
	"database/sql"

	"github.com/igorsfreitas/devbook-api/src/models"
)

// Posts represents a post repository
type Posts struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *Posts {
	return &Posts{db}
}

// Create creates a new post
func (repository Posts) Create(post models.Post) (uint64, error) {
	statement, err := repository.db.Prepare("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	err = statement.QueryRow(post.Title, post.Content, post.AuthorID).Scan(&post.ID)
	if err != nil {
		return 0, err
	}

	return uint64(post.ID), nil
}

// GetPost returns a single post by id
func (repository Posts) GetPost(postID uint64) (models.Post, error) {
	linha, err := repository.db.Query(
		`
			SELECT p.*, u.nick FROM posts p
			INNER JOIN users u ON u.id = p.author_id
			WHERE p.id = $1
		`,
		postID,
	)
	if err != nil {
		return models.Post{}, err
	}

	defer linha.Close()

	var post models.Post
	if linha.Next() {
		if err = linha.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

// GetPosts returns all posts from the user and the users that he follows
func (repository Posts) GetPosts(userID uint64) ([]models.Post, error) {
	linhas, err := repository.db.Query(
		`
			SELECT DISTINCT p.*, u.nick FROM posts p
			INNER JOIN users u ON u.id = p.author_id
			INNER JOIN followers f ON p.author_id = f.user_id
			WHERE u.id = $1 OR f.follower_id = $1
			ORDER BY 1 DESC
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var posts []models.Post
	for linhas.Next() {
		var post models.Post
		if err = linhas.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// Update updates a post
func (repository Posts) Update(postID uint64, post models.Post) error {
	statement, err := repository.db.Prepare("UPDATE posts SET title = $1, content = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(post.Title, post.Content, postID); err != nil {
		return err
	}

	return nil
}

// Delete deletes a post
func (repository Posts) Delete(postID uint64) error {
	statement, err := repository.db.Prepare("DELETE FROM posts WHERE id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(postID); err != nil {
		return err
	}

	return nil
}
