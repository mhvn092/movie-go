package staff

import "time"

type Staff struct {
	Id          int       `json:"id"            db:"id"`
	FirstName   string    `json:"first_name"    db:"first_name"    validate:"required, is_string"`
	LastName    string    `json:"last_name"     db:"last_name"     validate:"required, is_string"`
	Bio         string    `json:"bio"           db:"bio"           validate:"required, is_string"`
	StaffTypeId int       `json:"staff_type_id" db:"staff_type_id" validate:"required, is_int"`
	BirthDate   string    `json:"birth_date"    db:"birth_date"    validate:"required, is_date_string"`
	createdAt   time.Time `                     db:"created_at"`
	updatedAt   time.Time `                     db:"updated_at"`
}
