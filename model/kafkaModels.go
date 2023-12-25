package model

import (
	"fmt"
	"strings"
)

type NotificationType string

const (
	NOTIFICATION_TYPE_EMAIL NotificationType = "EMAIL"
)

type NotificationEvent struct {
	Recipients       []string
	Subject          string
	Body             string
	IsHtml           bool
	NotificationType NotificationType
}

func (event *NotificationEvent) String() string {
	return fmt.Sprintf("NotificationEvent: { 'recipients': %s,  subject: %s, isHtml: %t, type: %s}", strings.Join(event.Recipients, ", "), event.Subject, event.IsHtml, event.NotificationType)
}
