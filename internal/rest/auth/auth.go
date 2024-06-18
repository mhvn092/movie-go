package auth

import (
	"github.com/mhvn092/movie-go/internal/util"
	"net/http"
)

func AuthMux() http.Handler {
	authMux := http.NewServeMux()
	authMux.Handle("/signup", util.Post(http.HandlerFunc(signup)))

	return http.StripPrefix("/api/v1/auth", authMux)
}

func signup(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("You All Signed Up"))
}
