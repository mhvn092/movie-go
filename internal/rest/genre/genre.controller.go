package genre

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/mhvn092/movie-go/internal/models/genre"
	"github.com/mhvn092/movie-go/internal/util"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func getAll(w http.ResponseWriter, req *http.Request) {
	params, ok := util.GetPaginationParam(req)
	if !ok {
		return
	}

	res, nextCursor, err := models.GetAllGenresPaginated(db, params)
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
