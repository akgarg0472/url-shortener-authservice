package router

import (
	"github.com/go-chi/chi"

	Handlers "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
)

func DiscoveryRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/info", func(r chi.Router) {
		r.Get("/", Handlers.DiscoveryInfoHandler)
	})

	router.Route("/health", func(r chi.Router) {
		r.Get("/", Handlers.DiscoveryHealthHandler)
	})

	return router
}
