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
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.LoginRequestBodyValidator)
		r.Post("/", AuthHandler.LoginHandler)
	})

	router.Route("/signup", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.SignupRequestBodyValidator)
		r.Post("/", AuthHandler.SignupHandler)
	})

	router.Route("/validate-token", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.VerifyTokenRequestBodyValidator)
		r.Post("/", AuthHandler.VerifyTokenHandler)
	})

	router.Route("/logout", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.LogoutRequestBodyValidator)
		r.Post("/", AuthHandler.LogoutHandler)
	})

	router.Route("/forgot-password", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.ForgotPasswordRequestBodyValidator)
		r.Post("/", AuthHandler.ForgotPasswordHandler)
	})

	router.Route("/verify-reset-password", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Get("/", AuthHandler.VerifyResetPasswordHandler)
	})

	router.Route("/reset-password", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.ResetPasswordRequestBodyValidator)
		r.Post("/", AuthHandler.ResetPasswordHandler)
	})

	router.Route("/verify-admin", func(r chi.Router) {
		r.Use(Middlewares.AddRequestIdHeader)
		r.Use(Middlewares.ValidateRequestJSONContentType)
		r.Use(Middlewares.VerifyAdminRequestBodyHandler)
		r.Post("/", AuthHandler.VerifyAdminHandler)
	})

	return router
}
