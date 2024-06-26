package auth

import (
	"github.com/mhvn092/movie-go/pkg/router"
	"net/http"
)

func Router() *router.Router {
	r := router.NewRouter()
	r.Post("/signup", signup)

	return r
}

func signup(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("You All Signed Up"))
}
