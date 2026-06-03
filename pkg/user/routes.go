package user

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the user table
func Routes(config *config.Config) chi.Router {

	// Init Router
	userConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", userConfig.PostHandler)
	router.Get("/{id}", userConfig.GetByIdHandler)
	router.Get("/", userConfig.GetAllHandler)
	router.Patch("/{id}", userConfig.UpdateHandler)
	router.Delete("/{id}", userConfig.DeleteHandler)

	return router
}
