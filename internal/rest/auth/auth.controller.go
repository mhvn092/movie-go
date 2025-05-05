package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	user "github.com/mhvn092/movie-go/internal/models/user"
	validator "github.com/mhvn092/movie-go/pkg/Validator"
	"github.com/mhvn092/movie-go/pkg/env"
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

	if isAdmin {
		payload.Role = user.UserRole.ADMIN
	}

	if err := payload.RegisterUser(db); err != nil {
		if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "this user already exist", http.StatusConflict)
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
		fmt.Println("inside the validation")
		return
	}

	user, err := payload.CheckUser(db)
	if err != nil {
		if err.Error() == strconv.Itoa(http.StatusNotFound) {
			exception.HttpError(err, w, "this user was not found", http.StatusNotFound)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}

	err = user.ComparePasswords(payload.Password)
	if err != nil {
		exception.HttpError(err, w, "email or password is incorrect", http.StatusNotFound)
		return
	}

	token, err := createToken(user)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Write([]byte(token))
}

func createToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.Id,
			"email": user.Email,
			"role":  user.Role,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	secretKey := []byte(env.GetEnv(env.JWT_SECRET_KEY))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
