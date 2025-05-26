package authhandler

import (
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/internal/domain/user"
	validator "github.com/mhvn092/movie-go/pkg/Validator"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func signupAdmin(w http.ResponseWriter, req *http.Request) {
	signup(w, req, true)
}

func singnupUser(w http.ResponseWriter, req *http.Request) {
	signup(w, req, false)
}

func signup(w http.ResponseWriter, req *http.Request, isAdmin bool) {
	var payload user.User

	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	if err := service.Register(&payload, isAdmin); err != nil {
		if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "user already exists", http.StatusConflict)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}

	w.Write([]byte("You All Signed Up"))
}

func login(w http.ResponseWriter, req *http.Request) {
	var payload user.LoginDto

	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	u, err := service.Login(&payload)
	if err != nil {
		exception.HttpError(err, w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := service.GenerateToken(u)
	if err != nil {
		exception.DefaultInternalHttpError(w)

		return
	}

	w.Write([]byte(token))
}
