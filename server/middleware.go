package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/harbor-xyz/coding-project/contract"
	"github.com/harbor-xyz/coding-project/controller"
)

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
