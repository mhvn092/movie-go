package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/mhvn092/movie-go/internal/rest/middleware"
	"github.com/mhvn092/movie-go/internal/util"
	"github.com/mhvn092/movie-go/pkg/exception"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

// NewRouter initializes a new Router
func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

// ServeHTTP implements the http.Handler interface for Router
func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(res, req)
}

func (r *Router) Serve(url string) error {
	return http.ListenAndServe(url, r.mux)
}

func (r *Router) Use(middleware middleware.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// applyMiddlewares applies all registered middleware to a handler
func (r *Router) applyMiddlewares(handler http.Handler) http.Handler {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) handleMethod(
	path, method string,
	handlerFunc http.HandlerFunc,
	routeMiddlewares ...middleware.Middleware,
) {
	// Wrap the handlerFunc to enforce HTTP method
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, req)
	})

	// Apply route-level middleware (like JWT)
	handler := http.Handler(baseHandler)
	for _, m := range routeMiddlewares {
		handler = m(handler)
	}

	// Apply global middleware (like logging, CORS, etc.)
	handler = r.applyMiddlewares(handler)

	// Register with mux (this line is correct)
	r.mux.Handle(path, handler)
}

func (r *Router) GetWithPagination(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	wrapped := func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		limitStr := query.Get("limit")
		cursorStr := query.Get("cursor_id")

		var (
			limit  uint64 = 20
			cursor uint64 = 0
			err    error
		)

		if limitStr != "" {
			if limit, err = strconv.ParseUint(limitStr, 10, 64); err != nil || limit == 0 {
				exception.HttpError(
					errors.New("Invalid parameter"),
					w,
					"Invalid parameter",
					http.StatusBadRequest,
				)
				return
			}
		}

		if cursorStr != "" {
			if cursor, err = strconv.ParseUint(cursorStr, 10, 64); err != nil {
				exception.HttpError(
					errors.New("Invalid parameter"),
					w,
					"Invalid parameter",
					http.StatusBadRequest,
				)
				return
			}
		}

		param := util.PaginationParam{Limit: limit, CursorID: cursor}
		ctx := context.WithValue(req.Context(), util.PaginationKey, param)
		req = req.WithContext(ctx)

		handlerFunc(w, req)
	}

	r.handleMethod(path, http.MethodGet, wrapped, middlewares...)
}

// Post is a custom method to handle POST requests
func (r *Router) Post(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	r.handleMethod(path, http.MethodPost, handlerFunc, middlewares...)
}

func (r *Router) Get(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	r.handleMethod(path, http.MethodGet, handlerFunc, middlewares...)
}

func (r *Router) Put(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	r.handleMethod(path, http.MethodPut, handlerFunc, middlewares...)
}

func (r *Router) Patch(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	r.handleMethod(path, http.MethodPatch, handlerFunc, middlewares...)
}

func (r *Router) Delete(
	path string,
	handlerFunc http.HandlerFunc,
	middlewares ...middleware.Middleware,
) {
	r.handleMethod(path, http.MethodDelete, handlerFunc, middlewares...)
}

func (r *Router) AddSubRoute(path string, subRouter *Router) {
	subRouter.middlewares = append(r.middlewares, subRouter.middlewares...)
	// Trim any trailing slash from the path
	cleanPath := strings.TrimSuffix(path, "/")
	handler := r.applyMiddlewares(subRouter)

	r.mux.Handle(path, http.StripPrefix(cleanPath, handler))
}
