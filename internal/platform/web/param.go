package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/pkg/exception"
)

func GetIdFromParam(req *http.Request, w http.ResponseWriter) int {
	idString := req.PathValue("id")
	if idString == "" {
		exception.HttpError(
			errors.New("No Id Provided"),
			w,
			"No Id Provided",
			http.StatusBadRequest,
		)
		return 0
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return 0
	}

	return id
}
