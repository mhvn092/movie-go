package user

import (
	"strings"
	"time"

	"github.com/mhvn092/movie-go/internal/platform/security"
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

// PrepareCreate Prepare user for register
func (u *User) prepareToCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	hashedPassword, err := security.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword

	if u.Role == "" {
		u.Role = UserRole.NORMAL
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
