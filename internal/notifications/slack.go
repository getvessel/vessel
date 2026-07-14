package notifications

import (
	"fmt"

	"vessl.dev/vessl/internal/models"
)

func SendSlackNotification(webhookURL string, event *models.NotificationEvent) error {
	return postJSON(webhookURL, map[string]string{
		"text": fmt.Sprintf("*%s*\n%s\n%s", event.Title, event.Message, event.URL),
	})
}
