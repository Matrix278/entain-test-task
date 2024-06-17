package model

import "errors"

type Error struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func ErrUserNotFound() error {
	return errors.New("user not found")
}

func ErrTransactionAlreadyExists() error {
	return errors.New("transaction already exists")
}

func ErrInsufficientBalance() error {
	return errors.New("insufficient balance")
}
