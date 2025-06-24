package staffhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mhvn092/movie-go/internal/domain/staff"
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

func getSearchResults(w http.ResponseWriter, req *http.Request) {
	searchTerm := req.URL.Query().Get("term")
	if searchTerm == "" {
		return
	}

	res, err := service.GetSearchResults(searchTerm)
	if err != nil {
		fmt.Println("error is : ", err)
		exception.DefaultInternalHttpError(w)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Write([]byte(response))
}

func getDetail(w http.ResponseWriter, req *http.Request) {
	id := web.GetIdFromParam(req, w)
	if id == 0 {
		return
	}

	res, err := service.GetDetail(id)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		exception.DefaultInternalHttpError(w)
		return
	}

	w.Write([]byte(response))
}

func insert(w http.ResponseWriter, req *http.Request) {
	var payload staff.Staff
	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	staffId, err := service.Insert(&payload)
	if err != nil {
		println("err", err.Error())
		if err.Error() == strconv.Itoa(http.StatusNotFound) {
			exception.HttpError(err, w, "staff not found", http.StatusNotFound)
		} else if err.Error() == strconv.Itoa(http.StatusConflict) {
			exception.HttpError(err, w, "staff type already exists", http.StatusConflict)
		} else {
			exception.DefaultInternalHttpError(w)
		}
		return
	}

	w.Write([]byte(strconv.Itoa(staffId)))
}

func edit(w http.ResponseWriter, req *http.Request) {
	id := web.GetIdFromParam(req, w)
	if id == 0 {
		return
	}

	var payload staff.Staff
	if validator.JsonBodyHasErrors(req, w, &payload) {
		return
	}

	if err := service.Edit(id, &payload); err != nil {
		if err.Error() == strconv.Itoa(http.StatusNotFound) {
			exception.HttpError(err, w, "staff not found", http.StatusNotFound)
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
