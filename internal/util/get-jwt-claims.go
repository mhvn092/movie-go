package util

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/internal/models/user"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type UserClaims struct {
	Id    int               `json:"id"`
	Email string            `json:"email"`
	Role  user.UserRoleType `json:"role"`
	jwt.RegisteredClaims
}

// ClaimsFromContext extracts claims from the request context
func ClaimsFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(jwt.Claims)
	if !ok {
		return nil, false
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	return mapClaims, ok
}
