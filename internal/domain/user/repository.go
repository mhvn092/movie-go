package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/mhvn092/movie-go/internal/platform/repository"
)

type UserRepository struct {
	*repository.BaseRepository
}

func NewUserRepository(base *repository.BaseRepository) *UserRepository {
	return &UserRepository{BaseRepository: base}
}

func (r *UserRepository) isUserAlreadyRegisted(email string) error {
	var id int
	err := r.DB.QueryRow(
		context.Background(),
		"select id from person.users where email = $1",
		email,
	).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	return errors.New(strconv.Itoa(http.StatusConflict))
}

func (r *UserRepository) registerUser(u *User) error {
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

func (r *UserRepository) checkUser(login *LoginDto) (*User, error) {
	var user User

	err := r.DB.QueryRow(
		context.Background(),
		"select id, password, email, role from person.users where email = $1",
		login.Email,
	).Scan(&user.Id, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(strconv.Itoa(http.StatusNotFound))
		}
		return nil, err
	}

	return &user, nil
}
