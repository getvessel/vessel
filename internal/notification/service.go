package notification

import "context"

type Dispatcher interface {
	Send(event *NotificationEvent) error
}

type Service struct {
	repo       Repository
	dispatcher Dispatcher
}

func NewService(repo Repository, dispatcher Dispatcher) *Service {
	return &Service{repo: repo, dispatcher: dispatcher}
}

func (s *Service) GetIntegration(ctx context.Context) (*NotificationIntegration, error) {
	return s.repo.GetIntegration(ctx)
}

func (s *Service) SaveIntegration(ctx context.Context, n *NotificationIntegration) error {
	return s.repo.SaveIntegration(ctx, n)
}

func (s *Service) GetProjectPref(ctx context.Context, projectID string) (*ProjectNotificationPref, error) {
	return s.repo.GetProjectPref(ctx, projectID)
}

func (s *Service) SaveProjectPref(ctx context.Context, pref *ProjectNotificationPref) error {
	return s.repo.SaveProjectPref(ctx, pref)
}

func (s *Service) SendTest(channel, projectID string) error {
	event := &NotificationEvent{
		Title:     "Vessel Notification Test",
		Message:   "If you are reading this, your " + channel + " notification channel is configured correctly!",
		Level:     "success",
		ProjectID: projectID,
		URL:       "http://localhost:8080/settings/notifications",
	}
	return s.dispatcher.Send(event)
}
