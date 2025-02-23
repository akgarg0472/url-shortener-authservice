package router

import (
	"github.com/go-chi/chi"

	"github.com/akgarg0472/urlshortener-auth-service/internal/handler"
	"github.com/akgarg0472/urlshortener-auth-service/internal/middleware"
)

func AuthRouterV1() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/login", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.LoginRequestBodyValidator)
		r.Post("/", handler.LoginHandler)
	})

	router.Route("/signup", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.SignupRequestBodyValidator)
		r.Post("/", handler.SignupHandler)
	})

	router.Route("/validate-token", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.VerifyTokenRequestBodyValidator)
		r.Post("/", handler.VerifyTokenHandler)
	})

	router.Route("/logout", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.LogoutRequestBodyValidator)
		r.Post("/", handler.LogoutHandler)
	})

	router.Route("/forgot-password", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.ForgotPasswordRequestBodyValidator)
		r.Post("/", handler.ForgotPasswordHandler)
	})

	router.Route("/verify-reset-password", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Get("/", handler.VerifyResetPasswordHandler)
	})

	router.Route("/reset-password", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.ResetPasswordRequestBodyValidator)
		r.Post("/", handler.ResetPasswordHandler)
	})

	router.Route("/verify-admin", func(r chi.Router) {
		r.Use(middleware.AddRequestIdHeader)
		r.Use(middleware.ValidateRequestJSONContentType)
		r.Use(middleware.VerifyAdminRequestBodyHandler)
		r.Post("/", handler.VerifyAdminHandler)
	})

	return router
}
