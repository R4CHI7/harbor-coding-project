package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/harbor-xyz/coding-project/controller"
	"github.com/harbor-xyz/coding-project/database"
	"github.com/harbor-xyz/coding-project/repository"
	"github.com/harbor-xyz/coding-project/service"
)

func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	db := database.Get()
	userRepository := repository.NewUser(db)
	userAvailabilityRepository := repository.NewUserAvailability(db)
	eventRepository := repository.NewEvent(db)
	slotRepository := repository.NewSlot(db)

	userController := controller.NewUser(service.NewUser(userRepository, userAvailabilityRepository))
	eventController := controller.NewEvent(service.NewEvent(eventRepository, slotRepository))
	slotController := controller.NewSlot(service.NewSlot(slotRepository, userAvailabilityRepository))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(userIDContext)
			r.Post("/availability", userController.SetAvailability)
			r.Get("/availability", userController.GetAvailability)
			r.Get("/availability_overlap", userController.GetAvailabilityOverlap)
			r.Route("/events", func(r chi.Router) {
				r.Post("/", eventController.Create)
				r.Get("/", eventController.GetAll)
			})
			r.Route("/slots", func(r chi.Router) {
				r.Post("/", slotController.Create)
				r.Get("/", slotController.GetAll)
				r.Delete("/{slotID}", slotController.Delete)
			})
		})
	})

	return r
}
