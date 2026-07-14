package services

import (
	"context"
	"errors"
	"os"

	"vessl.dev/vessl/internal/models"
)

type NotificationDispatcher interface {
	Send(event *models.NotificationEvent) error
}

type NotificationService struct {
	dispatcher NotificationDispatcher
}

func NewNotificationService(dispatcher NotificationDispatcher) *NotificationService {
	return &NotificationService{dispatcher: dispatcher}
}

func (s *NotificationService) TestGlobalNotification(ctx context.Context, provider string) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher unavailable")
	}
	dashboardURL := os.Getenv("VESSL_DASHBOARD_URL")
	return s.dispatcher.Send(&models.NotificationEvent{
		EventType: "test_global_" + provider,
		Title:     "Global Test Notification from Vessl",
		Message:   "If you see this, your global integration is working correctly!",
		Level:     "info",
		URL:       dashboardURL + "/settings/notifications",
	})
}
