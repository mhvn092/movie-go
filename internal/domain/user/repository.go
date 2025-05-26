package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/internal/platform/repository"
)

type UserRepository struct {
	*repository.BaseRepository
}

func NewUserRepository(base *repository.BaseRepository) *UserRepository {
	return &UserRepository{BaseRepository: base}
}

func (r *UserRepository) isUserAlreadyRegisted(email string) error {
	rows, err := r.DB.Query(
		context.Background(),
		"select id from person.users where email = $1",
		email,
	)
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return err
		}
	} else {
		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	}

	return errors.New(strconv.Itoa(http.StatusConflict))
}

func (r *UserRepository) RegisterUser(u *User) error {
	if err := r.isUserAlreadyRegisted(u.Email); err != nil {
		return err
	}
	err := u.prepareToCreate()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(
		context.Background(),
		"Insert into person.users (first_name, last_name, email, password, role, phone_number, created_at, updated_at) values($1, $2, $3,$4, $5,$6,$7,$8)",
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.Role,
		u.PhoneNumber,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CheckUser(login *LoginDto) (*User, error) {
	rows, err := r.DB.Query(
		context.Background(),
		"select id, password, email, role from person.users where email = $1",
		login.Email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user User
	if rows.Next() {
		if err := rows.Scan(&user.Id, &user.Password, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		return &user, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New(strconv.Itoa(http.StatusNotFound))
}
