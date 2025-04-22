package exception

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mhvn092/movie-go/pkg/env"
)

func ErrorExit(e error, message string) {
	if e != nil {
		if env.GetEnv(env.ENVIROMENT) == "development" {
			err := fmt.Sprintf("%s: the real error(%s)", message, e.Error())
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func HttpError(e error, w http.ResponseWriter, message string, code int) {
	if e != nil {
		if env.GetEnv(env.ENVIROMENT) == "development" {
			err := fmt.Sprintf("%s: the real error(%s)", message, e.Error())
			fmt.Println(err)
		}
		http.Error(w, message, code)
	}
}

func DefaultInternalHttpError(w http.ResponseWriter) {
	HttpError(
		errors.New("Some Unexpected Error Happened, Please Try again later"),
		w,
		"Some Unexpected Error Happened, Please Try again later",
		http.StatusInternalServerError,
	)
}
