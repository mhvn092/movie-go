package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/pkg/env"
)

type UserTokenData struct {
	ID    int
	Email string
	Role  string
}

func CreateToken(data UserTokenData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    data.ID,
		"email": data.Email,
		"role":  data.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := []byte(env.GetEnv(env.JWT_SECRET_KEY))
	return token.SignedString(secretKey)
}
