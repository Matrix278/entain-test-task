package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type Transaction struct {
	ID         strfmt.UUID4 `json:"id"`
	UserID     strfmt.UUID4 `json:"user_id"`
	Amount     float64      `json:"amount"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  *time.Time   `json:"updated_at,omitempty"`
	CanceledAt *time.Time   `json:"canceled_at,omitempty"`
}
