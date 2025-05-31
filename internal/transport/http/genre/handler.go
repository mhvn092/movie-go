package genrehandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/internal/domain/genre"
	"github.com/mhvn092/movie-go/internal/platform/web"
	validator "github.com/mhvn092/movie-go/pkg/Validator"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func getAll(w http.ResponseWriter, req *http.Request) {
	params, ok := web.GetPaginationParam(req)
	if !ok {
		return
	}

	res, nextCursor, err := service.GetAllPaginated(params)
	if err != nil {
		fmt.Println(err)
		exception.DefaultInternalHttpError(w)
	}

	w.Header().Add("X-Next-Cursor", fmt.Sprintf("%d", nextCursor))

	response, err := json.Marshal(res)
	if err != nil {
		exception.DefaultInternalHttpError(w)
	}

	w.Write([]byte(response))
}

func insert(w http.ResponseWriter, req *http.Request) {
	var payload genre.Genre
	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	genreId, err := service.Insert(&payload)
	if err != nil {
		if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "genre already exists", http.StatusConflict)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}

	w.Write([]byte(strconv.Itoa(genreId)))
}

func edit(w http.ResponseWriter, req *http.Request) {
}

func delete(w http.ResponseWriter, req *http.Request) {
}
