package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRoleType string

var UserRole = struct {
	ADMIN  UserRoleType
	NORMAL UserRoleType
}{
	ADMIN:  "admin",
	NORMAL: "normal",
}

// User full model
type User struct {
	Id          int          `db:"id"`
	FirstName   string       `db:"first_name"   json:"first_name"   validate:"is_string"`
	LastName    string       `db:"last_name"    json:"last_name"    validate:"is_string"`
	Email       string       `db:"email"        json:"email"        validate:"required, is_string, is_email"`
	Password    string       `db:"password"     json:"password"     validate:"required, is_string, is_strong_password,min_len=10"`
	Role        UserRoleType `db:"role"`
	PhoneNumber string       `db:"phone_number" json:"phone_number" validate:"required, is_string, is_phone_number"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}

type LoginDto struct {
	Email    string `json:"email"    validate:"required, is_string, is_email"`
	Password string `json:"password" validate:"required, is_string"`
}

// HashPassword Hash user password with bcrypt
func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePasswords Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// PrepareCreate Prepare user for register
func (u *User) prepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.hashPassword(); err != nil {
		return err
	}

	if u.Role == "" {
		u.Role = UserRole.NORMAL
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) isUserAlreadyRegisted(db *pgxpool.Pool) error {
	rows, err := db.Query(
		context.Background(),
		"select id from person.users where email = $1",
		u.Email,
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

func (u *User) RegisterUser(db *pgxpool.Pool) error {
	if err := u.isUserAlreadyRegisted(db); err != nil {
		return err
	}
	err := u.prepareCreate()
	if err != nil {
		return err
	}

	_, err = db.Exec(
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

func (login *LoginDto) CheckUser(db *pgxpool.Pool) (*User, error) {
	rows, err := db.Query(
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
