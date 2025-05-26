package models

import "time"

type Staff struct {
	Id          int       `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Bio         string    `json:"bio" db:"bio"`
	StaffTypeId int       `json:"staff_type_id" db:"staff_type_id"`
	BirthDate   time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
