package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/internal/rest/middleware"
	"github.com/mhvn092/movie-go/pkg/router"
)

var db *pgxpool.Pool

func Router() *router.Router {
	db = config.GetDbPool()
	r := router.NewRouter()
	r.Post("/signup", singnupUser)
	r.Post("/login", login)
	r.Post("/add-operator", signupAdmin, middleware.IsAdminAuthorized())
	return r
}
