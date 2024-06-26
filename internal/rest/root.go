package root

import (
	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/internal/rest/auth"
	"github.com/mhvn092/movie-go/pkg/exception"
	"net/http"
)

func InitializeRoutes() {
	r := config.AppConfig.Mux
	globalPrefix := "/api/v1/"
	r.Get("/", rootHandler)
	r.AddSubRoute(globalPrefix+"auth/", auth.Router())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	exception.HttpError(err, w, "some error exists", 500)
}
