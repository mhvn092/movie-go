package auth

import (
	"github.com/mhvn092/movie-go/pkg/router"
)

func Router() *router.Router {
	r := router.NewRouter()
	r.Post("/signup", signup)
	return r
}
