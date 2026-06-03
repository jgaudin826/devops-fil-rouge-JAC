package field

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the field table
func Routes(config *config.Config) chi.Router {

	// Init Router
	fieldConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", fieldConfig.PostHandler)
	router.Get("/{id}", fieldConfig.GetByIdHandler)
	router.Get("/", fieldConfig.GetAllHandler)
	router.Patch("/{id}", fieldConfig.UpdateHandler)
	router.Delete("/{id}", fieldConfig.DeleteHandler)
	router.Get("/{id}/media", fieldConfig.GetMediaByFieldHandler)

	return router
}
