package main

import (
	"log"
	"mediadex/config"
	"mediadex/pkg/authentication"
	"mediadex/pkg/collection"
	"mediadex/pkg/field"
	"mediadex/pkg/media"
	"mediadex/pkg/mediaField"
	"mediadex/pkg/tag"
	"mediadex/pkg/user"
	"net/http"

	_ "mediadex/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Initialize CORS
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)

	// Swagger endpoint
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Authentication endpoint
	router.Mount("/api/v1/auth", authentication.Routes(configuration))

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(authentication.AuthMiddleware(configuration.JWTSecret))
		r.Mount("/collections", collection.Routes(configuration))
		r.Mount("/fields", field.Routes(configuration))
		r.Mount("/media", media.Routes(configuration))
		r.Mount("/mediaFields", mediaField.Routes(configuration))
		r.Mount("/tags", tag.Routes(configuration))
		r.Mount("/users", user.Routes(configuration))
	})

	return router
}

// @title MediaDex API
// @version 1.0
// @description This is the MediaDex API
// @host
// @BasePath /api/v1
// @securityDefinitions.apikey	BearerAuth
// @in				header
// @name			Authorization
func main() {
	// Configuration Initialization
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error:", err)
	}

	// Routes Initialization
	router := Routes(configuration)

	// Serve the API
	log.Println("Server running on http://localhost:" + configuration.Constants.Port)
	log.Println("Swagger UI available at http://localhost:" + configuration.Constants.Port + "/swagger/index.html")
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.Port, router))
}
