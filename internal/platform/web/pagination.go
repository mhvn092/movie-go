package web

import "net/http"

type PaginationParam struct {
	Limit    uint64
	CursorID uint64
}

const PaginationKey string = "pagination_param"

func GetPaginationParam(r *http.Request) (PaginationParam, bool) {
	p, ok := r.Context().Value(PaginationKey).(PaginationParam)
	return p, ok
}
