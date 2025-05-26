package genrehandler

import (
	"github.com/mhvn092/movie-go/internal/domain/genre"
	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/middleware"
	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/pkg/router"
)

var service *genre.GenreService

func initialize() {
	db := config.GetDbPool()
	genreRepo := genre.NewGenreRepository(&repository.BaseRepository{DB: db})
	service = genre.NewGenreService(genreRepo)
}

func Router() *router.Router {
	initialize()
	r := router.NewRouter()

	r.GetWithPagination("/all", getAll)
	r.Post("/create", insert, middleware.AuthAdmin)
	r.Put("/update/:id", edit, middleware.AuthAdmin)
	r.Delete("/delete/:id", delete, middleware.AuthAdmin)
	return r
}
