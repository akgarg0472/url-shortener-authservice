package router

import (
	"github.com/go-chi/chi"

	"github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	"github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func OAuthRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/providers", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Get("/", handler.GetOAuthProvidersHandler)
	})

	router.Route("/callbacks", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.OAuthCallbackRequestBodyValidator)
		r.Post("/", handler.OAuthCallbackHandler)
	})

	return router
}
