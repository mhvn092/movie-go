package stafftype

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/internal/platform/web"
)

type StaffTypeRepository struct {
	*repository.BaseRepository
}

func NewStaffTypeRepository(base *repository.BaseRepository) *StaffTypeRepository {
	return &StaffTypeRepository{BaseRepository: base}
}

func (r *StaffTypeRepository) getAllStaffTypesPaginated(
	params web.PaginationParam,
) (res []StaffType, nextCursor int, err error) {
	rows, err := r.DB.Query(
		context.Background(),
		"select id, title from staff.staff_type where id >= $1 limit $2",
		params.CursorID,
		params.Limit,
	)

	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	res = []StaffType{}

	for rows.Next() {
		var item StaffType
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

func (r *StaffTypeRepository) checkIfExists(query string, args ...interface{}) (bool, error) {
	var staffTypeId int
	err := r.DB.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&staffTypeId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *StaffTypeRepository) checkIfExistsByTitle(title string) (bool, error) {
	return r.checkIfExists("select id from staff.staff_type where title = $1", title)
}

func (r *StaffTypeRepository) checkIfExistsById(id int) (bool, error) {
	return r.checkIfExists("select id from staff.staff_type where id = $1", id)
}

func (r *StaffTypeRepository) checkIfExistsByNameExcludingId(id int, title string) (bool, error) {
	return r.checkIfExists(
		"select id from staff.staff_type where title = $1 and id <> $2",
		title,
		id,
	)
}

func (r *StaffTypeRepository) insert(staffType *StaffType) (int, error) {
	exists, err := r.checkIfExistsByTitle(staffType.Title)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errors.New(strconv.Itoa(http.StatusConflict))
	}

	var staffTypeId int

	rows, err := r.DB.Query(
		context.Background(),
		"insert into staff.staff_type (title) values ($1) returning id",
		staffType.Title,
	)

	defer rows.Close()

	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&staffTypeId)
		if err != nil {
			return 0, err
		}
	}

	return staffTypeId, nil
}

func (r *StaffTypeRepository) edit(id int, staffType *StaffType) error {
	exists, err := r.checkIfExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}

	exists, err = r.checkIfExistsByNameExcludingId(id, staffType.Title)
	if err != nil {
		return err
	}

	if exists {
		return errors.New(strconv.Itoa(http.StatusConflict))
	}

	cmdTag, err := r.DB.Exec(
		context.Background(),
		"update staff.staff_type set title = $1 where id = $2",
		staffType.Title,
		id,
	)
	if err != nil {
		return err
	}

	println("cmdTag", cmdTag.RowsAffected())
	if cmdTag.RowsAffected() == 0 {
		return errors.New("Could not update")
	}

	return nil
}

func (r *StaffTypeRepository) delete(id int) error {
	exists, err := r.checkIfExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(strconv.Itoa(http.StatusNotFound))
	}
	cmdTag, err := r.DB.Exec(
		context.Background(),
		"delete from staff.staff_type where id = $1",
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
