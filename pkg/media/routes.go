package media

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the media table
func Routes(config *config.Config) chi.Router {

	// Init Router
	mediaConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", mediaConfig.PostHandler)
	router.Get("/{id}", mediaConfig.GetByIdHandler)
	router.Get("/", mediaConfig.GetAllHandler)
	router.Patch("/{id}", mediaConfig.UpdateHandler)
	router.Delete("/{id}", mediaConfig.DeleteHandler)
	router.Get("/{id}/tags", mediaConfig.GetTagsByMediaHandler)
	router.Get("/{id}/fields", mediaConfig.GetFieldsByMediaHandler)

	return router
}
