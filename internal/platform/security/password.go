package security

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePasswords compares a bcrypt hashed password with its possible plain-text equivalent.
func ComparePasswords(hashedPassword, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plain))
}
