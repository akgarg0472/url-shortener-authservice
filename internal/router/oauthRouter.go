package router

import (
	"github.com/go-chi/chi"

	OAuthHandler "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
)

func OAuthRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/clients", func(r chi.Router) {
		r.Get("/", OAuthHandler.GetOAuthClients)
	})

	return router
}
