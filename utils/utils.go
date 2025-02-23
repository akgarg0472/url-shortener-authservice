package utils

import (
	"strings"
)

func GenerateForgotPasswordTokenRedirectUrl(email string, token string) string {
	frontendBaseUrl := GetEnvVariable("FRONTEND_BASE_DOMAIN", "http://127.0.0.1:3000/")
	resetPasswordPageUrl := GetEnvVariable("FRONTEND_RESET_PASSWORD_PAGE_URL", "reset-password")
	return EnsureTrailingSlash(frontendBaseUrl) + resetPasswordPageUrl + "?token=" + token + "&email=" + email
}

func GenerateForgotPasswordLink(email string, forgotPasswordToken string) string {
	backendBaseUrl := GetEnvVariable("BACKEND_BASE_DOMAIN", "http://localhost:8765/")
	backendResetPasswordUrl := GetEnvVariable("BACKEND_RESET_PASSWORD_URL", "auth/v1/reset-password")
	return EnsureTrailingSlash(backendBaseUrl) + backendResetPasswordUrl + "?email=" + email + "&token=" + forgotPasswordToken
}

func GetStringOrNil(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func GetInt64OrNil(i *int64) int64 {
	if i != nil {
		return *i
	}

	return -1
}

func EnsureTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}
