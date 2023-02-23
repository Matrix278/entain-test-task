package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type Transaction struct {
	ID         strfmt.UUID4 `db:"id"`
	UserID     strfmt.UUID4 `db:"user_id"`
	Amount     float64      `db:"amount"`
	CreatedAt  time.Time    `db:"created_at"`
	CanceledAt *time.Time   `db:"canceled_at"`
}
