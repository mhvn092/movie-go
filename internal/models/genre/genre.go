package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mhvn092/movie-go/internal/util"
)

type Genre struct {
	Id        int       `json:"id"                   db:"id"`
	Title     string    `json:"title"                db:"title"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func GetAllGenresPaginated(
	db *pgxpool.Pool,
	params util.PaginationParam,
) (res []Genre, nextCursor int, err error) {
	rows, err := db.Query(
		context.Background(),
		"select id, title from movie.genre where id >= $1 limit $2",
		params.CursorID,
		params.Limit,
	)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	res = []Genre{}

	for rows.Next() {
		var item Genre
		err = rows.Scan(
			&item.Id,
			&item.Title,
		)
		if err != nil {
			return
		}
		res = append(res, item)
	}

	if len(res) > 0 {
		nextCursor = res[len(res)-1].Id + 1
	}

	return
}
