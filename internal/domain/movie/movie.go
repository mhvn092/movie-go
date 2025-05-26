package models

import "time"

type Movie struct {
	Id             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	DirectorId     int       `json:"director_id" db:"director_id"`
	GenreId        int       `json:"genre_id" db:"genre_id"`
	ProductionYear int       `json:"production_year" db:"production_year"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type MovieStaff struct {
	MovieId     int `json:"movie_id" db:"movie_id"`
	StaffId     int `json:"staff_id" db:"staff_id"`
	StaffTypeId int `json:"staff_type_id" db:"staff_type_id"`
}
