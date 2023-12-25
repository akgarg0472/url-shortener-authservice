package utils

import "strings"

func GenerateForgotPasswordTokenRedirectUrl(email string, token string) string {
	frontendBaseUrl := GetEnvVariable("FRONTEND_BASE_DOMAIN", "http://127.0.0.1:3000/")
	resetPasswordPageUrl := GetEnvVariable("FRONTEND_RESET_PASSWORD_PAGE_URL", "reset-password")
	return frontendBaseUrl + resetPasswordPageUrl + "?token=" + token + "&email=" + email
}

func GenerateResetPasswordLink(email string, forgotPasswordToken string) string {
	backendBaseUrl := GetEnvVariable("BACKEND_BASE_DOMAIN", "http://localhost:8765/")
	backendResetPasswordUrl := GetEnvVariable("BACKEND_RESET_PASSWORD_URL", "auth/v1/reset-password")
	return backendBaseUrl + backendResetPasswordUrl + "?email=" + email + "&token=" + forgotPasswordToken
}

func GetFormattedName(firstName string, lastName string) string {
	if len(strings.TrimSpace(firstName)) == 0 {
		return lastName
	}

	return firstName + " " + lastName
}
