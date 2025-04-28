package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	user "github.com/mhvn092/movie-go/internal/models/user"
	validator "github.com/mhvn092/movie-go/pkg/Validator"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func signup(w http.ResponseWriter, req *http.Request) {
	var payload user.User

	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	if payload.RegisterUser(w, db) {
		w.Write([]byte("You All Signed Up"))
	}
}

func login(w http.ResponseWriter, req *http.Request) {
	var payload user.LoginDto

	if validator.JsonBodyHasErrors(req, w, &payload) {
		fmt.Println("inside the validation")
		return
	}

	user, ok := payload.CheckUser(w, db)
	if !ok {
		return
	}

	err := user.ComparePasswords(payload.Password)
	if err != nil {
		exception.HttpError(err, w, "email or password is incorrect", http.StatusNotFound)
		return
	}

	token, err := createToken(user.Email, user.Id)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Write([]byte(token))
}

func createToken(email string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    id,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	secretKey := []byte(env.GetEnv(env.JWT_SECRET_KEY))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
