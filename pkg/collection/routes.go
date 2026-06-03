package collection

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

// Routes the collection table
func Routes(config *config.Config) chi.Router {

	// Init Router
	collectionConfig := New(config)
	router := chi.NewRouter()

	// Routes
	router.Post("/", collectionConfig.PostHandler)
	router.Get("/{id}", collectionConfig.GetByIdHandler)
	router.Get("/", collectionConfig.GetAllHandler)
	router.Patch("/{id}", collectionConfig.UpdateHandler)
	router.Delete("/{id}", collectionConfig.DeleteHandler)

	return router
}
