package stafftypehandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	stafftype "github.com/mhvn092/movie-go/internal/domain/staff-type"
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
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Header().Add("X-Next-Cursor", fmt.Sprintf("%d", nextCursor))

	response, err := json.Marshal(res)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Write([]byte(response))
}

func insert(w http.ResponseWriter, req *http.Request) {
	var payload stafftype.StaffType
	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	staffTypeId, err := service.Insert(&payload)
	if err != nil {
		if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "staff type already exists", http.StatusConflict)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}

	w.Write([]byte(strconv.Itoa(staffTypeId)))
}

func edit(w http.ResponseWriter, req *http.Request) {
	id := web.GetIdFromParam(req, w)
	if id == 0 {
		return
	}

	var payload stafftype.StaffType
	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	if err := service.Edit(id, &payload); err != nil {
		if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "staff type already exists", http.StatusConflict)
		} else if err.Error() == strconv.Itoa(http.StatusNotFound) {
			exception.HttpError(err, w, "staff type not found", http.StatusNotFound)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))
}

func delete(w http.ResponseWriter, req *http.Request) {
	id := web.GetIdFromParam(req, w)
	if id == 0 {
		return
	}

	if err := service.Delete(id); err != nil {
		if err.Error() == strconv.Itoa(http.StatusNotFound) {
			exception.HttpError(err, w, "staff type not found", http.StatusNotFound)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Success"))
}
