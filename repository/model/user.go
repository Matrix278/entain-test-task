package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type User struct {
	ID        strfmt.UUID4 `db:"id"`
	Balance   float64      `db:"balance"`
	CreatedAt time.Time    `db:"created_at"`
}
