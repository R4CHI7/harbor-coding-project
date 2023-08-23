package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/harbor-xyz/coding-project/controller"
	"github.com/harbor-xyz/coding-project/database"
	"github.com/harbor-xyz/coding-project/repository"
	"github.com/harbor-xyz/coding-project/service"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	userController := controller.NewUser(service.NewUser(repository.NewUser(database.Get())))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
	})

	return r
}
