package staff

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

type StaffRepository struct {
	*repository.BaseRepository
}

func NewStaffRepository(base *repository.BaseRepository) *StaffRepository {
	return &StaffRepository{BaseRepository: base}
}

func (r *StaffRepository) getAllStaffPaginated(
	params web.PaginationParam,
) (res []StaffGetAllResponse, nextCursor int, err error) {
	rows, err := r.DB.Query(
		context.Background(),
		"select id, first_name, last_name from staff.staff where id >= $1 limit $2",
		params.CursorID,
		params.Limit,
	)

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	res = []StaffGetAllResponse{}

	for rows.Next() {
		var item StaffGetAllResponse
		err = rows.Scan(
			&item.Id,
			&item.FirstName,
			&item.LastName,
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

func (r *StaffRepository) getSearchResults(
	searchTerm string,
) (res []StaffGetAllResponse, err error) {
	terms := strings.Fields(searchTerm)
	queryTerm := strings.Join(terms, " <-> ") + ":*"
	rows, err := r.DB.Query(
		context.Background(),
		`SELECT id, first_name, last_name
    FROM staff.staff
    WHERE search_vector @@ to_tsquery('simple', $1)
    LIMIT 10;`,
		queryTerm,
	)

	defer rows.Close()

	if err != nil {
		fmt.Println("the error is the in the query itself")
		return nil, err
	}

	res = []StaffGetAllResponse{}

	for rows.Next() {
		var item StaffGetAllResponse
		err = rows.Scan(
			&item.Id,
			&item.FirstName,
			&item.LastName,
		)
		if err != nil {
			return
		}
		res = append(res, item)
	}

	return
}

func (r *StaffRepository) checkIfExists(id int) (bool, error) {
	var staffId int
	err := r.DB.QueryRow(
		context.Background(),
		"select id from staff.staff where id = $1",
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

func (r *StaffRepository) checkCountOfExistingIds(ids []int) (bool, error) {
	currentCount := len(ids)
	if currentCount == 0 {
		return false, nil
	}

	placeholders := make([]string, currentCount)
	args := make([]interface{}, currentCount)

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		"select count(id) from staff.staff where id in (%s)",
		strings.Join(placeholders, ","),
	)

	var existingCount int
	err := r.DB.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&existingCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return existingCount == currentCount, nil
}

func (r *StaffRepository) getDetail(id int) (StaffGetDetailResponse, error) {
	var staff StaffGetDetailResponse
	err := r.DB.QueryRow(
		context.Background(),
		"SELECT s.id as id , first_name, last_name, bio, birth_date, staff_type_id, st.title as staff_type_title from staff.staff s inner join staff.staff_type st on staff_type_id = st.id where s.id = $1",
		id,
	).Scan(&staff.Id, &staff.FirstName, &staff.LastName, &staff.Bio, &staff.BirthDate, &staff.StaffTypeId, &staff.StaffTypeTitle)
	if err != nil {
		if err == pgx.ErrNoRows {
			return staff, errors.New(strconv.Itoa(http.StatusNotFound))
		}
		return staff, err
	}
	return staff, nil
}

func (r *StaffRepository) insert(staff *Staff) (int, error) {
	var staffId int

	rows, err := r.DB.Query(
		context.Background(),
		"insert into staff.staff (first_name, last_name, bio, birth_date, staff_type_id) values ($1,$2,$3,$4,$5) returning id",
		staff.FirstName,
		staff.LastName,
		staff.Bio,
		staff.BirthDate,
		staff.StaffTypeId,
	)

	defer rows.Close()

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&staffId)
		if err != nil {
			return 0, err
		}
	}

	return staffId, nil
}

func (r *StaffRepository) edit(id int, staff *Staff) error {
	exists, err := r.checkIfExists(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	cmdTag, err := r.DB.Exec(
		context.Background(),
		"update staff.staff set first_name = $1, last_name = $2, bio = $3, birth_date = $4, staff_type_id =$5, updated_at = $6 where id = $7",
		staff.FirstName,
		staff.LastName,
		staff.Bio,
		staff.BirthDate,
		staff.StaffTypeId,
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

func (r *StaffRepository) delete(id int) error {
	exists, err := r.checkIfExists(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	cmdTag, err := r.DB.Exec(
		context.Background(),
		"delete from staff.staff where id = $1",
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
