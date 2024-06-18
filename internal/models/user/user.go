package models

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
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
	Id          int           `json:"id" db:"id"`
	FirstName   string        `json:"first_name" db:"first_name"`
	LastName    string        `json:"last_name" db:"last_name"`
	Email       string        `json:"email,omitempty" db:"email"`
	Password    string        `json:"password,omitempty" db:"password"`
	Role        *UserRoleType `json:"role,omitempty" db:"role"`
	PhoneNumber *string       `json:"phone_number,omitempty" db:"phone_number"`
	CreatedAt   time.Time     `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty" db:"updated_at"`
}

// HashPassword Hash user password with bcrypt
func (u *User) HashPassword() error {
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

// SanitizePassword Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// PrepareCreate Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role == nil {
		*u.Role = UserRole.NORMAL
	}
	return nil
}
