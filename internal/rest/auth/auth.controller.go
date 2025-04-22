package auth

import (
	"net/http"

	config "github.com/mhvn092/movie-go/internal"
	models "github.com/mhvn092/movie-go/internal/models/user"
	validator "github.com/mhvn092/movie-go/pkg/Validator"
)

func signup(res http.ResponseWriter, req *http.Request) {
	var payload models.User

	if validator.JsonBodyHasErrors(req, res, &payload) {
		return
	}

	db := config.GetDbPool()

	if payload.RegisterUser(res, db) {
		res.Write([]byte("You All Signed Up"))
	}
}
