package authhandler

import (
	"github.com/mhvn092/movie-go/internal/domain/user"
	"github.com/mhvn092/movie-go/internal/platform/config"
	"github.com/mhvn092/movie-go/internal/platform/middleware"
	"github.com/mhvn092/movie-go/internal/platform/repository"
	"github.com/mhvn092/movie-go/pkg/router"
)

var service *user.UserService

func initialize() {
	db := config.GetDbPool()
	userRepo := user.NewUserRepository(&repository.BaseRepository{DB: db})
	service = user.NewUserService(userRepo)
}

func Router() *router.Router {
	initialize()
	r := router.NewRouter()
	r.Post("/signup", singnupUser)
	r.Post("/login", login)
	r.Post("/add-operator", signupAdmin, middleware.AuthAdmin)
	return r
}
