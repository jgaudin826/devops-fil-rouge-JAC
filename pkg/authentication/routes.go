package authentication

import (
	"mediadex/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) chi.Router {
	UserConfig := New(configuration)
	router := chi.NewRouter()
	router.Post("/login", UserConfig.LoginHandler)
	router.Post("/refresh", UserConfig.RefreshHandler)
	router.Post("/register", UserConfig.RegisterHandler)
	return router
}
