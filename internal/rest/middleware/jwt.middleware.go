package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/mhvn092/movie-go/internal/models/user"
	"github.com/mhvn092/movie-go/internal/util"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func IsAdminAuthorized() Middleware {
	return authorized(true)
}

func IsUserAuthorized() Middleware {
	return authorized(false)
}

func authorized(checkAdmin bool) Middleware {
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
				&util.UserClaims{},
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
				claims, ok := token.Claims.(*util.UserClaims)
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

			ctx := context.WithValue(r.Context(), util.ClaimsKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
