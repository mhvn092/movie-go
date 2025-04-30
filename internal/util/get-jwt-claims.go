package util

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/internal/rest/middleware"
)

// ClaimsFromContext extracts claims from the request context
func ClaimsFromContext(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(jwt.Claims)
	if !ok {
		return nil, false
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	return mapClaims, ok
}
