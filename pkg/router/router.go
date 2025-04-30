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

func (r *Router) handleMethod(
	path, method string,
	handlerFunc http.HandlerFunc,
	routeMiddlewares ...Middleware,
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

// Post is a custom method to handle POST requests
func (r *Router) Post(path string, handlerFunc http.HandlerFunc, middlewares ...Middleware) {
	r.handleMethod(path, http.MethodPost, handlerFunc, middlewares...)
}

func (r *Router) Get(path string, handlerFunc http.HandlerFunc, middlewares ...Middleware) {
	r.handleMethod(path, http.MethodGet, handlerFunc, middlewares...)
}

func (r *Router) Put(path string, handlerFunc http.HandlerFunc, middlewares ...Middleware) {
	r.handleMethod(path, http.MethodPut, handlerFunc, middlewares...)
}

func (r *Router) Patch(path string, handlerFunc http.HandlerFunc, middlewares ...Middleware) {
	r.handleMethod(path, http.MethodPatch, handlerFunc, middlewares...)
}

func (r *Router) Delete(path string, handlerFunc http.HandlerFunc, middlewares ...Middleware) {
	r.handleMethod(path, http.MethodDelete, handlerFunc, middlewares...)
}

func (r *Router) AddSubRoute(path string, subRouter *Router) {
	subRouter.middlewares = append(r.middlewares, subRouter.middlewares...)
	// Trim any trailing slash from the path
	cleanPath := strings.TrimSuffix(path, "/")
	handler := r.applyMiddlewares(subRouter)

	r.mux.Handle(path, http.StripPrefix(cleanPath, handler))
}
