package model

type Error struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}
