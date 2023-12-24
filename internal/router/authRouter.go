package router

import (
	"github.com/go-chi/chi"

	AuthHandler "github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	Middlewares "github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func AuthRouterV1() *chi.Mux {
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

	router.Route("/validate-token", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestContentType)
		r.Use(Middlewares.VerifyTokenRequestBodyValidator)
		r.Post("/", AuthHandler.VerifyToken)
	})

	router.Route("/logout", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestContentType)
		r.Use(Middlewares.LogoutRequestBodyValidator)
		r.Post("/", AuthHandler.Logout)
	})

	router.Route("/forgot-password", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestContentType)
		r.Use(Middlewares.ForgotPasswordRequestBodyValidator)
		r.Post("/", AuthHandler.ForgotPassword)
	})

	router.Route("/reset-password", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Get("/", AuthHandler.ResetPassword)
	})

	return router
}
