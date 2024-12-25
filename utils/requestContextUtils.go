package utils

type contextKey string

// RequestContextKeys holds the context key constants
var RequestContextKeys = struct {
	LoginRequestKey          contextKey
	SignupRequestKey         contextKey
	LogoutRequestKey         contextKey
	ValidateTokenRequestKey  contextKey
	ForgotPasswordRequestKey contextKey
	OAuthCallbackRequestKey  contextKey
	ResetPasswordRequestKey  contextKey
}{
	LoginRequestKey:          "loginRequest",
	SignupRequestKey:         "signupRequest",
	LogoutRequestKey:         "logoutRequest",
	ValidateTokenRequestKey:  "validateTokenRequest",
	ForgotPasswordRequestKey: "forgotPasswordRequest",
	OAuthCallbackRequestKey:  "oAuthCallbackRequest",
	ResetPasswordRequestKey:  "resetPasswordRequest",
}
