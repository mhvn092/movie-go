package genrehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mhvn092/movie-go/internal/platform/web"
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
}

func edit(w http.ResponseWriter, req *http.Request) {
}

func delete(w http.ResponseWriter, req *http.Request) {
}
