package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/internal/domain/user"
	"github.com/mhvn092/movie-go/internal/platform/security"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func isAdminAuthorized() Middleware {
	return authorized(true)
}

func isUserAuthorized() Middleware {
	return authorized(false)
}

func authorized(checkAdmin bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := env.GetEnv(env.JWT_SECRET_KEY)

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				exception.HttpError(
					errors.New("Missing Authorization header"),
					w, "Missing Authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(
				tokenStr,
				&security.UserClaims{},
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method")
					}
					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				exception.HttpError(err, w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*security.UserClaims)
			if !ok {
				exception.HttpError(
					errors.New("Invalid claims"),
					w,
					"Invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			if checkAdmin && claims.Role != string(user.UserRole.ADMIN) {
				exception.HttpError(errors.New("Forbidden"), w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), security.ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
