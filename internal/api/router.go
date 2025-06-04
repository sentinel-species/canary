package api

import (
	"canary/internal/database"
	"canary/internal/services/dependency"
	"canary/internal/store"
	"net/http"
	"time"

	"canary/internal/api/handlers"
	"canary/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db database.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	healthHandler := handlers.NewHealthHandler()

	dependencyStore := store.NewDependencyStore(db)
	dependencyService := dependency.NewService(*dependencyStore)

	dependencyHandler := handlers.NewDependencyHandler(dependencyService)

	r.Get("/health", healthHandler.HealthCheck())
	r.Post("/api/v1/dependency/resolve", dependencyHandler.ResolveDependency())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		handlers.RespondJSON(w, http.StatusNotFound, types.StatusResponse{
			Status:  "error",
			Message: "endpoint not found",
		})
	})

	return r
}
