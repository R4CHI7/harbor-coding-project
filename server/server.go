package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/harbor-xyz/coding-project/contract"
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
	userController := controller.NewUser(service.NewUser(repository.NewUser(db), repository.NewUserAvailability(db)))
	eventController := controller.NewEvent(service.NewEvent(repository.NewEvent(db), repository.NewSlot(db)))
	slotController := controller.NewSlot(service.NewSlot(repository.NewSlot(db), repository.NewUserAvailability(db)))

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
			})
		})
	})

	return r
}

func userIDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Render(w, r, contract.ErrorRenderer(errors.New("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			render.Render(w, r, contract.ErrorRenderer(errors.New("invalid user ID")))
		}
		ctx := context.WithValue(r.Context(), controller.ContextUserIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
