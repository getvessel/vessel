package dispatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type DispatcherService struct {
	repo repositories.NotificationRepository
}

func NewDispatcherService(repo repositories.NotificationRepository) *DispatcherService {
	return &DispatcherService{repo: repo}
}

func (d *DispatcherService) Dispatch(event *models.NotificationEvent) {
	go func() {
		if err := d.Send(event); err != nil {
			log.Printf("[Dispatcher] Failed to dispatch event '%s': %v", event.Title, err)
		}
	}()
}

func (d *DispatcherService) Send(event *models.NotificationEvent) error {
	if event.TeamID == "" {
		return fmt.Errorf("TeamID is required for dispatch")
	}

	channels, err := d.repo.ListChannelsByTeam(context.Background(), event.TeamID)
	if err != nil {
		return fmt.Errorf("failed to list channels for team %s: %w", event.TeamID, err)
	}

	for _, c := range channels {
		if !c.IsEnabled {
			continue
		}

		// Check if event type matches
		var events []string
		if len(c.Events) > 0 {
			if err := json.Unmarshal(c.Events, &events); err == nil {
				matches := false
				for _, e := range events {
					if e == event.EventType || e == "*" {
						matches = true
						break
					}
				}
				if !matches && len(events) > 0 {
					continue
				}
			}
		}

		switch c.Provider {
		case "slack":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = d.sendWebhook(cfg.WebhookURL, event)
			}
		case "discord":
			var cfg struct {
				WebhookURL string `json:"webhookUrl"`
			}
			if json.Unmarshal(c.Config, &cfg) == nil && cfg.WebhookURL != "" {
				_ = d.sendWebhook(cfg.WebhookURL, event)
			}
		case "smtp":
			// ... simplified smtp logic here if needed
		}
	}
	return nil
}

func (d *DispatcherService) sendWebhook(webhookURL string, event *models.NotificationEvent) error {
	payload := map[string]string{
		"content": fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	return err
}
