package model

import (
	"fmt"
	"strings"

	enums "github.com/akgarg0472/urlshortener-auth-service/constants"
)

type NotificationEvent struct {
	Recipients       []string
	Subject          string
	Body             string
	IsHtml           bool
	NotificationType enums.NotificationType
}

func (event *NotificationEvent) String() string {
	return fmt.Sprintf("NotificationEvent: { 'recipients': %s,  subject: %s, isHtml: %t, type: %s}", strings.Join(event.Recipients, ", "), event.Subject, event.IsHtml, event.NotificationType)
}

type UserRegisteredEvent struct {
	UserId string `json:"user_id"`
}
