package genre

import (
	"context"

	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type GenreRepository struct {
	*repository.BaseRepository
}

func NewGenreRepository(base *repository.BaseRepository) *GenreRepository {
	return &GenreRepository{BaseRepository: base}
}

func (r *GenreRepository) GetAllGenresPaginated(
	params web.PaginationParam,
) (res []Genre, nextCursor int, err error) {
	rows, err := r.DB.Query(
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
