package root

import (
	"fmt"
	"net/http"

	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/internal/rest/auth"
	"github.com/mhvn092/movie-go/internal/rest/genre"
	"github.com/mhvn092/movie-go/internal/rest/middleware"
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
	r.Use(middleware.RequestLogger())

	r.Get("/", rootHandler)

	r.AddSubRoute(globalPrefix+"auth/", auth.Router())
	r.AddSubRoute(globalPrefix+"genre/", genre.Router())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	exception.HttpError(err, w, "some error exists", 500)
}
