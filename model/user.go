package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type User struct {
	ID        strfmt.UUID4 `json:"id"`
	Balance   float64      `json:"balance"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at,omitempty"`
}
