package tag

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the tag table
func Routes(config *config.Config) chi.Router {

	// Init Router
	tagConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", tagConfig.PostHandler)
	router.Get("/{id}", tagConfig.GetByIdHandler)
	router.Get("/", tagConfig.GetAllHandler)
	router.Patch("/{id}", tagConfig.UpdateHandler)
	router.Delete("/{id}", tagConfig.DeleteHandler)
	router.Get("/{id}/media", tagConfig.GetMediaByTagHandler)

	return router
}
