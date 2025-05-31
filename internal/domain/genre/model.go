package genre

import (
	"time"
)

type Genre struct {
	Id        int       `json:"id"                   db:"id"`
	Title     string    `json:"title"                db:"title"      validate:"required, is_string"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
