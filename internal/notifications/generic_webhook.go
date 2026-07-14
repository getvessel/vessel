package notifications

import (
	"vessl.dev/vessl/internal/models"
)

func SendGenericWebhook(webhookURL string, event *models.NotificationEvent) error {
	return postJSON(webhookURL, event)
}
