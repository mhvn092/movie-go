package root

import (
	config "github.com/mhvn092/movie-go/internal"
	"github.com/mhvn092/movie-go/internal/rest/auth"
	"github.com/mhvn092/movie-go/internal/util"
	"net/http"
)

func InitializeRoutes() {
	mux := config.AppConfig.Mux
	mux.HandleFunc("/", rootHandler)
	mux.Handle("/api/v1/auth/", auth.AuthMux())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	util.HttpError(err, w, "")
}
