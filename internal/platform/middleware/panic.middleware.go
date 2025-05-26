package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/mhvn092/movie-go/pkg/exception"
)

// RecoverPanic handles panics and returns a 500 Internal Server Error.
func recoverPanic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Log the panic with stack trace
					log.Printf("[Panic] %v\n%s", rec, debug.Stack())

					// Return a 500 Internal Server Error
					exception.HttpError(
						fmt.Errorf("internal server error"),
						w,
						"Internal Server Error",
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
