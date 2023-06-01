package models

import (
	"errors"

	"github.com/igorsfreitas/devbook-api/src/commons/encrypt"
)

// Password represents a password update model
type Password struct {
	NewPassword     string `json:"newPassword"`
	CurrentPassword string `json:"currentPassword"`
}

// Prepare prepares the password model
func (password *Password) Prepare() error {
	if err := password.validate(); err != nil {
		return err
	}

	if err := password.format(); err != nil {
		return err
	}

	return nil
}

func (password *Password) validate() error {
	if password.NewPassword == "" {
		return errors.New("a nova senha não pode estar em branco")
	}

	if password.CurrentPassword == "" {
		return errors.New("a senha atual não pode estar em branco")
	}

	return nil
}

func (password *Password) format() error {

	hash, err := encrypt.HashPassword(password.NewPassword)
	if err != nil {
		return err
	}

	password.NewPassword = string(hash)

	return nil
}
