package router

import (
	"net/http"
	"strings"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

type Middleware func(http.Handler) http.Handler

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

func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// applyMiddlewares applies all registered middleware to a handler
func (r *Router) applyMiddlewares(handler http.Handler) http.Handler {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) handleMethod(path, method string, handlerFunc http.HandlerFunc) {
	wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, req)
	})

	handler := http.Handler(wrappedHandler)
	// Apply middlewares to the wrapped handler
	handler = r.applyMiddlewares(handler)

	// Register the handler with the path
	r.mux.Handle(path, handler)
}

// Post is a custom method to handle POST requests
func (r *Router) Post(path string, handlerFunc http.HandlerFunc) {
	r.handleMethod(path, http.MethodPost, handlerFunc)
}

func (r *Router) Get(path string, handlerFunc http.HandlerFunc) {
	r.handleMethod(path, http.MethodGet, handlerFunc)
}

func (r *Router) Put(path string, handlerFunc http.HandlerFunc) {
	r.handleMethod(path, http.MethodPut, handlerFunc)
}

func (r *Router) Patch(path string, handlerFunc http.HandlerFunc) {
	r.handleMethod(path, http.MethodPatch, handlerFunc)
}

func (r *Router) Delete(path string, handlerFunc http.HandlerFunc) {
	r.handleMethod(path, http.MethodDelete, handlerFunc)
}

func (r *Router) AddSubRoute(path string, subRouter *Router) {
	subRouter.middlewares = append(r.middlewares, subRouter.middlewares...)
	// Trim any trailing slash from the path
	cleanPath := strings.TrimSuffix(path, "/")
	handler := r.applyMiddlewares(subRouter)

	r.mux.Handle(path, http.StripPrefix(cleanPath, handler))
}
