package model

import "errors"

type Error struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func ErrUserNotFound() error {
	return errors.New("user not found")
}
