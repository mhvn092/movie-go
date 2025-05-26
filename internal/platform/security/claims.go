package security

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type UserClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func ClaimsFromContext(r *http.Request) (*UserClaims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(*UserClaims)
	return claims, ok
}
