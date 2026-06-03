package mediaField

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the mediaField table
func Routes(config *config.Config) chi.Router {

	// Init Router
	mediaFieldConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", mediaFieldConfig.PostHandler)
	router.Get("/{fieldID}/{mediaID}", mediaFieldConfig.GetByIdHandler)
	router.Get("/", mediaFieldConfig.GetAllHandler)
	router.Patch("/{fieldID}/{mediaID}", mediaFieldConfig.UpdateHandler)
	router.Delete("/{fieldID}/{mediaID}", mediaFieldConfig.DeleteHandler)

	return router
}
