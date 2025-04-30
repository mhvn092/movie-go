package root

import (
	"net/http"

	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/internal/rest/auth"
	"github.com/mhvn092/movie-go/internal/rest/middleware"
	"github.com/mhvn092/movie-go/pkg/exception"
)

func InitializeRoutes() {
	r := config.GetRouter()
	globalPrefix := "/api/v1/"
	r.Use(middleware.RequestLogger())

	r.Get("/", rootHandler)

	r.AddSubRoute(globalPrefix+"auth/", auth.Router())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	exception.HttpError(err, w, "some error exists", 500)
}
