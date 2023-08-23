package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

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

	userController := controller.NewUser(service.NewUser(repository.NewUser(database.Get())))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.Create)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(userIDContext)
			r.Post("/availability", userController.SetAvailability)
		})
	})

	return r
}

func userIDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Render(w, r, contract.ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			render.Render(w, r, contract.ErrorRenderer(fmt.Errorf("invalid user ID")))
		}
		ctx := context.WithValue(r.Context(), controller.ContextUserIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
