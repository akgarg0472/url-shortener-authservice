package router

import (
	"github.com/go-chi/chi"

	"github.com/akgarg0472/urlshortener-auth-service/internal/handler"
)

func PingRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/ping", func(r chi.Router) {
		r.Get("/", handler.PingHandler)
	})

	return router
}
