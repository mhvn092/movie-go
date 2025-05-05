package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/internal/models/user"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
	"github.com/mhvn092/movie-go/pkg/router"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type UserClaims struct {
	Id    int               `json:"id"`
	Email string            `json:"email"`
	Role  user.UserRoleType `json:"role"`
	jwt.RegisteredClaims
}

func IsAdminAuthorized() router.Middleware {
	return authorized(true)
}

func IsUserAuthorized() router.Middleware {
	return authorized(false)
}

func authorized(checkAdmin bool) router.Middleware {
	secret := env.GetEnv(env.JWT_SECRET_KEY)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				exception.HttpError(
					errors.New("Missing Authorization header"),
					w,
					"Missing Authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(
				tokenStr,
				&UserClaims{},
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method")
					}
					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				exception.HttpError(
					err,
					w,
					"Invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			if checkAdmin {
				claims, ok := token.Claims.(*UserClaims)
				if !ok {
					exception.HttpError(
						errors.New("Forbidden"),
						w,
						"Forbidden",
						http.StatusForbidden)
					return
				}

				if claims.Role != user.UserRole.ADMIN {
					exception.HttpError(
						errors.New("Forbidden"),
						w,
						"Forbidden",
						http.StatusForbidden)
					return
				}

			}

			ctx := context.WithValue(r.Context(), ClaimsKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
