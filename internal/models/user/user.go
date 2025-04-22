package models

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/mhvn092/movie-go/pkg/exception"
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

func (u *User) isUserAlreadyRegisted(w http.ResponseWriter, db *pgxpool.Pool) bool {
	rows, err := db.Query(
		context.Background(),
		"select id from person.users where email = $1",
		u.Email,
	)
	defer rows.Close()

	if err != nil {
		exception.DefaultQueryFailedHttpError(w, "user is already registered")
		return true
	}

	if rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			exception.DefaultQueryFailedHttpError(w, "user is already registered")
		}
	} else {
		if err := rows.Err(); err != nil {
			exception.DefaultQueryFailedHttpError(w, "user is already registered")
		}
		return false
	}
	return true
}

func (u *User) RegisterUser(w http.ResponseWriter, db *pgxpool.Pool) bool {
	if u.isUserAlreadyRegisted(w, db) {
		return false
	}
	err := u.prepareCreate()
	if err != nil {
		exception.DefaultQueryFailedHttpError(
			w,
			"some Error in preparing to create the user happened",
		)
		return false
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
		exception.DefaultQueryFailedHttpError(w, "Some Error in Registering user happened")
		return false
	}

	return true
}
