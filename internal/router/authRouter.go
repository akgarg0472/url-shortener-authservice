package router

import (
	"github.com/go-chi/chi"

	authHandler "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	middlewares "github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func AuthRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/login", func(r chi.Router) {
		r.Use(middlewares.AddRequestIdHeader)
		r.Use(middlewares.ValidateRequestContentType)
		r.Use(middlewares.LoginRequestBodyValidator)
		r.Post("/", authHandler.Login)
	})

	router.Route("/signup", func(r chi.Router) {
		r.Use(middlewares.AddRequestIdHeader)
		r.Use(middlewares.ValidateRequestContentType)
		r.Use(middlewares.SignupRequestBodyValidator)
		r.Post("/", authHandler.Signup)
	})

	return router
}
