package models

import (
	"errors"
	"strings"
	"time"
)

// Post model represents a user Post in the database
type Post struct {
	ID         uint64    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorID   uint64    `json:"authorId,omitempty"`
	AuthorNick string    `json:"authorNick,omitempty"`
	Likes      uint64    `json:"likes"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

// Prepare prepares the post to be saved
func (post *Post) Prepare() error {
	if err := post.validate(); err != nil {
		return err
	}

	post.format()

	return nil
}

func (post *Post) validate() error {
	if post.Title == "" {
		return errors.New("the title is required")
	}

	if post.Content == "" {
		return errors.New("the content is required")
	}

	return nil
}

func (post *Post) format() {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)
}
