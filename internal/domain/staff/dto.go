package staff

import "time"

type StaffGetAllResponse struct {
	Id        int    `json:"id"         db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name"  db:"last_name"`
}

type StaffGetDetailResponse struct {
	Id             int       `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Bio            string    `json:"bio"`
	StaffTypeId    int       `json:"staff_type_id"`
	StaffTypeTitle string    `json:"staff_type_title"`
	BirthDate      time.Time `json:"birth_date"`
}
