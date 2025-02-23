package notification_service

import (
	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	kafka_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

func SendSignupSuccessEmail(requestId string, email string, name string) {
	if logger.IsInfoEnabled() {
		logger.Info("Pushing signup success email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", email),
		)
	}

	body := utils.GetSignupSuccessEmailBody(name)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Welcome Aboard! Start Enjoying Link Shortening Bliss 🚀🎉", body, true, constants.NotificationTypeEmail)

	kafka_service.GetInstance().PushNotificationEvent(requestId, *event)
}

func SendForgotPasswordEmail(
	requestId string,
	email string,
	name string,
	forgotPasswordUrl string,
) {
	if logger.IsInfoEnabled() {
		logger.Info("Pushing forgot password email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", email),
		)
	}

	body := utils.GenerateForgotPasswordEmailBody(email, name, forgotPasswordUrl)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Reset your UrlShortener password", body, true, constants.NotificationTypeEmail)

	kafka_service.GetInstance().PushNotificationEvent(requestId, *event)
}

func SendPasswordChangeSuccessEmail(requestId string, email string) {
	if logger.IsInfoEnabled() {
		logger.Info("Pushing password changed success email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", email),
		)
	}

	body := utils.GeneratePasswordChangeSuccessEmailBody(email)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Password changed successfully 🎉", body, true, constants.NotificationTypeEmail)

	kafka_service.GetInstance().PushNotificationEvent(requestId, *event)
}

func generateNotificationEvent(
	recipients []string,
	subject string,
	body string,
	html bool,
	notificationType constants.NotificationType,
) *model.NotificationEvent {
	return &model.NotificationEvent{
		Recipients:       recipients,
		Subject:          subject,
		Body:             body,
		IsHtml:           html,
		NotificationType: notificationType,
	}
}
