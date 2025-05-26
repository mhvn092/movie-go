package root

import (
	"fmt"
	"net/http"

	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/middleware"
	authhandler "github.com/mhvn092/movie-go/internal/transport/http/auth"
	genrehandler "github.com/mhvn092/movie-go/internal/transport/http/genre"
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
	globalPrefix := "/api/v1/"
	r.Use(middleware.Logger)
	r.Use(middleware.RecoverPanic)

	r.Get("/", rootHandler)

	r.AddSubRoute(globalPrefix+"auth/", authhandler.Router())
	r.AddSubRoute(globalPrefix+"genre/", genrehandler.Router())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	exception.HttpError(err, w, "some error exists", 500)
}
