package router

import (
	"github.com/go-chi/chi"

	OAuthHandler "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	Middlewares "github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func OAuthRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/providers", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Get("/", OAuthHandler.GetOAuthProvidersHandler)
	})

	router.Route("/callbacks", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.OAuthCallbackRequestBodyValidator)
		r.Post("/", OAuthHandler.OAuthCallbackHandler)
	})

	return router
}
