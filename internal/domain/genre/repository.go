package genre

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type GenreRepository struct {
	*repository.BaseRepository
}

func NewGenreRepository(base *repository.BaseRepository) *GenreRepository {
	return &GenreRepository{BaseRepository: base}
}

func (r *GenreRepository) getAllGenresPaginated(
	params web.PaginationParam,
) (res []Genre, nextCursor int, err error) {
	rows, err := r.DB.Query(
		context.Background(),
		"select id, title from movie.genre where id >= $1 limit $2",
		params.CursorID,
		params.Limit,
	)

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

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

func (r *GenreRepository) checkIfExists(query string, args ...interface{}) (bool, error) {
	var genreId int
	err := r.DB.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&genreId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *GenreRepository) checkIfExistsByTitle(title string) (bool, error) {
	return r.checkIfExists("select id from movie.genre where title = $1", title)
}

func (r *GenreRepository) checkIfExistsById(id int) (bool, error) {
	return r.checkIfExists("select id from movie.genre where id = $1", id)
}

func (r *GenreRepository) checkIfExistsByNameExcludingId(id int, title string) (bool, error) {
	return r.checkIfExists("select id from movie.genre where title = $1 and id <> $2", title, id)
}

func (r *GenreRepository) insert(genre *Genre) (int, error) {
	exists, err := r.checkIfExistsByTitle(genre.Title)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errors.New(strconv.Itoa(http.StatusConflict))
	}

	var genreId int

	rows, err := r.DB.Query(
		context.Background(),
		"insert into movie.genre (title) values ($1) returning id",
		genre.Title,
	)

	defer rows.Close()

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&genreId)
		if err != nil {
			return 0, err
		}
	}

	return genreId, nil
}

func (r *GenreRepository) edit(id int, genre *Genre) error {
	exists, err := r.checkIfExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	exists, err = r.checkIfExistsByNameExcludingId(id, genre.Title)
	if err != nil {
		return err
	}

	if exists {
		return errors.New(strconv.Itoa(http.StatusConflict))
	}

	cmdTag, err := r.DB.Exec(
		context.Background(),
		"update movie.genre set title = $1, updated_at = $2 where id = $3",
		genre.Title,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("Could not update")
	}

	return nil
}

func (r *GenreRepository) delete(id int) error {
	exists, err := r.checkIfExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}
	cmdTag, err := r.DB.Exec(
		context.Background(),
		"delete from movie.genre where id = $1",
		id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("Could not delete")
	}

	return nil
}
