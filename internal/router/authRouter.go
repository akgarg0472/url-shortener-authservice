package router

import (
	"github.com/go-chi/chi"

	AuthHandler "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	Middlewares "github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func AuthRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/login", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestContentType)
		r.Use(Middlewares.LoginRequestBodyValidator)
		r.Post("/", AuthHandler.Login)
	})

	router.Route("/signup", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestContentType)
		r.Use(Middlewares.SignupRequestBodyValidator)
		r.Post("/", AuthHandler.Signup)
	})

	return router
}
