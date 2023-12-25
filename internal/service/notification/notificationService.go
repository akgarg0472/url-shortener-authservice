package notification_service

import (
	kafkaService "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger = Logger.GetLogger("notificationService.go")
)

func SendSignupSuccessEmail(requestId string, email string, name string) {
	logger.Info("[{}] Sending signup success email to {}", requestId, email)

	body := utils.GetSignupSuccessEmailBody(email, name)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Welcome Aboard: Link Shortening Bliss! 🚀🎉", body, true, model.NOTIFICATION_TYPE_EMAIL)

	kafkaService.GetInstance().PushNotificationEvent(*event)
}

func SendForgotPasswordEmail(requestId string, email string, name string, forgotPasswordUrl string) bool {
	logger.Info("[{}] Sending forgot password email to {}", requestId, email)

	body := utils.GenerateForgotPasswordEmailBody(email, name, forgotPasswordUrl)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Reset your UrlShortener password", body, true, model.NOTIFICATION_TYPE_EMAIL)

	return kafkaService.GetInstance().PushNotificationEvent(*event)
}

func SendPasswordChangeSuccessEmail(requestId string, email string) {
	logger.Info("[{}] Sending password changed success email to {}", requestId, email)

	body := utils.GeneratePasswordChangeSuccessEmailBody(email)
	recipients := [1]string{email}
	event := generateNotificationEvent(recipients[:], "Password changed successfully 🎉", body, true, model.NOTIFICATION_TYPE_EMAIL)

	kafkaService.GetInstance().PushNotificationEvent(*event)
}

func generateNotificationEvent(recipients []string, subject string, body string, html bool, notificationType model.NotificationType) *model.NotificationEvent {
	return &model.NotificationEvent{
		Recipients:       recipients,
		Subject:          subject,
		Body:             body,
		IsHtml:           html,
		NotificationType: notificationType,
	}
}
