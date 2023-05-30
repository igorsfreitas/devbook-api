package models

import (
	"errors"
	"strings"
	"time"
)

// User representa um usuário utilizando a rede social
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Prepare valida e formata o usuário
func (user *User) Prepare() error {
	if err := user.validate(); err != nil {
		return err
	}

	user.format()
	return nil
}

func (user *User) validate() (err error) {
	if user.Name == "" {
		err = errors.New("o nome é obrigatório e não pode estar em branco")
	}
	if user.Nick == "" {
		err = errors.New("o nick é obrigatório e não pode estar em branco")
	}
	if user.Email == "" {
		err = errors.New("o email é obrigatório e não pode estar em branco")
	}
	if user.Password == "" {
		err = errors.New("a senha é obrigatória e não pode estar em branco")
	}
	return
}

func (user *User) format() {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)
}
