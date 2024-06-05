package models

import "time"

type Movie struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	DirectorId  int       `json:"director_id" db:"director_id"`
	GenreId     int       `json:"genre_id" db:"genre_id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
