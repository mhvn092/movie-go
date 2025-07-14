package root

import (
	"fmt"
	"net/http"

	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/middleware"
	authhandler "github.com/mhvn092/movie-go/internal/transport/http/auth"
	genrehandler "github.com/mhvn092/movie-go/internal/transport/http/genre"
	moviehandler "github.com/mhvn092/movie-go/internal/transport/http/movie"
	staffhandler "github.com/mhvn092/movie-go/internal/transport/http/staff"
	stafftypehandler "github.com/mhvn092/movie-go/internal/transport/http/staff-type"
	"github.com/mhvn092/movie-go/pkg/env"
	"github.com/mhvn092/movie-go/pkg/exception"
	"github.com/mhvn092/movie-go/pkg/router"
)

func CreateServer() (string, *router.Router) {
	r := router.NewRouter()

	host := env.GetEnv(env.HOST)
	port := env.GetEnv(env.PORT)

	url := host + ":" + port
	fmt.Println("Listening on " + url)

	return url, r
}

func InitializeRoutes() {
	r := config.GetRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RecoverPanic)

	r.Get("/", rootHandler)

	r.AddSubRoute(getSubRoute("auth"), authhandler.Router())
	r.AddSubRoute(getSubRoute("genre"), genrehandler.Router())
	r.AddSubRoute(getSubRoute("staff-type"), stafftypehandler.Router())
	r.AddSubRoute(getSubRoute("staff"), staffhandler.Router())
	r.AddSubRoute(getSubRoute("movie"), moviehandler.Router())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	exception.HttpError(err, w, "some error exists", 500)
}

func getSubRoute(subRoute string) string {
	return "/api/v1/" + subRoute + "/"
}
