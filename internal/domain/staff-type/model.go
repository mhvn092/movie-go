package stafftype

import "time"

type StaffType struct {
	Id        int       `json:"id"    db:"id"`
	Title     string    `json:"title" db:"title"      validate:"required, is_string"`
	createdAt time.Time `             db:"created_at"`
	updatedAt time.Time `             db:"updated_at"`
}
