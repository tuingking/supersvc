package mux

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/tuingking/supersvc/handler/api"
	xmiddleware "github.com/tuingking/supersvc/pkg/middleware"
)

func NewMux(h api.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(xmiddleware.Inbound)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	})
	r.Use(cors.Handler)

	r.Get("/", h.HandleRoot)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/configs", h.GetAppConfig)

		// user API
		r.Get("/users", h.GetUser)
		r.Post("/users", h.CreateUser)
	})

	return r
}
