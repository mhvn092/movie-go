package movie

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type MovieRepository struct {
	*repository.BaseRepository
}

func NewMovieRepository(base *repository.BaseRepository) *MovieRepository {
	return &MovieRepository{BaseRepository: base}
}

func (r *MovieRepository) getAllMoviePaginated(
	params web.PaginationParam,
) (res []MovieGetAllResponse, nextCursor int, err error) {
	rows, err := r.DB.Query(
		context.Background(),
		"select id, title, production_year from movie.movie where id >= $1 limit $2",
		params.CursorID,
		params.Limit,
	)

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	res = []MovieGetAllResponse{}

	for rows.Next() {
		var item MovieGetAllResponse
		err = rows.Scan(
			&item.Id,
			&item.Title,
			&item.ProductionYear,
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

func (r *MovieRepository) getSearchResults(
	searchTerm string,
) (res []MovieGetAllResponse, err error) {
	terms := strings.Fields(searchTerm)
	queryTerm := strings.Join(terms, " <-> ") + ":*"
	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id, title, production_year
    FROM movie.movie
    WHERE search_vector @@ to_tsquery('simple', $1)
    LIMIT 10;`,
		queryTerm,
	)

	defer rows.Close()

	if err != nil {
		fmt.Println("the error is the in the query itself")
		return nil, err
	}

	res = []MovieGetAllResponse{}

	for rows.Next() {
		var item MovieGetAllResponse
		err = rows.Scan(
			&item.Id,
			&item.Title,
			&item.ProductionYear,
		)
		if err != nil {
			return
		}
		res = append(res, item)
	}

	return
}

func (r *MovieRepository) checkIfExists(id int) (bool, error) {
	var staffId int
	err := r.DB.QueryRow(
		context.Background(),
		"select id from movie.movie where id = $1",
		id,
	).Scan(&staffId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *MovieRepository) getDetail(id int) (MovieGetDetailResponse, error) {
	var movie MovieGetDetailResponse
	query := `
        SELECT 
            m.id,
            m.title,
            m.production_year,
            m.director_id,
            CONCAT(s.first_name, ' ', s.last_name) AS director_name, 
            m.genre_id,
            g.title AS genre_name,
            m.description,
            staff_data.staff_id,
            staff_data.staff_name,
            staff_data.staff_type_id,
            staff_data.staff_type_title
        FROM movie.movie m
        JOIN staff.staff s ON m.director_id = s.id
        JOIN movie.genre g ON m.genre_id = g.id
        LEFT JOIN LATERAL (
            SELECT 
                ms.staff_id AS staff_id,
                CONCAT(st.first_name, ' ', st.last_name) AS staff_name,
                ms.staff_type_id AS staff_type_id,
                stt.title AS staff_type_title
            FROM movie.movie_staff ms
            JOIN staff.staff st ON ms.staff_id = st.id
            JOIN staff.staff_type stt ON ms.staff_type_id = stt.id
            WHERE ms.movie_id = m.id
        ) staff_data ON true
        WHERE m.id = $1`

	rows, err := r.DB.Query(context.Background(), query, id)
	if err != nil {
		return movie, fmt.Errorf("failed to query movie detail: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		if errors.Is(rows.Err(), pgx.ErrNoRows) || rows.Err() == nil {
			return movie, errors.New(strconv.Itoa(http.StatusNotFound))
		}
		return movie, fmt.Errorf("failed to get first row: %w", rows.Err())
	}

	var staff movieStaffGetDetailBaseRow
	err = rows.Scan(
		&movie.Id,
		&movie.Title,
		&movie.ProductionYear,
		&movie.DirectorId,
		&movie.DirectorName,
		&movie.GenreId,
		&movie.GenreName,
		&movie.Description,
		&staff.StaffId,
		&staff.StaffName,
		&staff.StaffTypeId,
		&staff.StaffTypeTitle,
	)
	if err != nil {
		return movie, fmt.Errorf("failed to scan first row: %w", err)
	}

	if staff.StaffId != nil {
		movie.Staffs = append(movie.Staffs, MovieStaffBaseResponse{
			StaffId:        *staff.StaffId,
			StaffName:      *staff.StaffName,
			StaffTypeId:    *staff.StaffTypeId,
			StaffTypeTitle: *staff.StaffTypeTitle,
		})
	}

	for rows.Next() {
		err := rows.Scan(
			new(int),
			new(string),
			new(int),
			new(int),
			new(string),
			new(int),
			new(string),
			new(string),
			&staff.StaffId,
			&staff.StaffName,
			&staff.StaffTypeId,
			&staff.StaffTypeTitle,
		)
		if err != nil {
			return movie, fmt.Errorf("failed to scan subsequent staff row: %w", err)
		}

		if staff.StaffId != nil {
			movie.Staffs = append(movie.Staffs, MovieStaffBaseResponse{
				StaffId:        *staff.StaffId,
				StaffName:      *staff.StaffName,
				StaffTypeId:    *staff.StaffTypeId,
				StaffTypeTitle: *staff.StaffTypeTitle,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return movie, fmt.Errorf("error after iterating through rows: %w", err)
	}

	return movie, nil
}

func (r *MovieRepository) Insert(payload *MovieUpsertPayload) (movieId int, err error) {
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	err = tx.QueryRow(
		ctx,
		"INSERT INTO movie.movie (title, description, production_year, director_id, genre_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		payload.Title,
		payload.Description,
		payload.ProductionYear,
		payload.DirectorId,
		payload.GenreId,
	).Scan(&movieId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert movie: %w", err)
	}

	staffQuery, args := getMovieStaffInsertQuery(movieId, payload)

	_, err = tx.Exec(ctx, staffQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert movie staff: %w", err)
	}

	return movieId, nil
}

func (r *MovieRepository) Edit(id int, payload *MovieUpsertPayload) (err error) {
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	cmdTag, err := tx.Exec(
		ctx,
		"UPDATE movie.movie SET title = $1, description = $2, production_year = $3, director_id = $4, genre_id = $5, updated_at = $6 WHERE id = $7",
		payload.Title,
		payload.Description,
		payload.ProductionYear,
		payload.DirectorId,
		payload.GenreId,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no movie found with id %d", id)
	}

	_, err = tx.Exec(ctx, "DELETE FROM movie.movie_staff WHERE movie_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete existing movie staff: %w", err)
	}

	staffQuery, args := getMovieStaffInsertQuery(id, payload)
	_, err = tx.Exec(ctx, staffQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert movie staff: %w", err)
	}

	return nil
}

func (r *MovieRepository) delete(id int) (err error) {
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return errors.New("could not start transaction")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	_, err = tx.Exec(ctx, "delete from movie.movie_staff where movie_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "delete from movie.movie where id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func getMovieStaffInsertQuery(
	id int,
	payload *MovieUpsertPayload,
) (query string, args []interface{}) {
	query = ""
	if len(payload.Staffs) > 0 {
		var baseStaffsQuery strings.Builder
		baseStaffsQuery.WriteString(
			"INSERT INTO movie.movie_staff (movie_id, staff_id, staff_type_id) VALUES ",
		)
		args = make([]interface{}, 0, len(payload.Staffs)*3)
		placeholders := make([]string, 0, len(payload.Staffs))

		for i, staff := range payload.Staffs {
			start := i*3 + 1
			placeholders = append(
				placeholders,
				fmt.Sprintf("($%d, $%d, $%d)", start, start+1, start+2),
			)
			args = append(args, id, staff.StaffId, staff.StaffTypeId)
		}

		baseStaffsQuery.WriteString(strings.Join(placeholders, ","))
		query = baseStaffsQuery.String()
	}
	return
}
