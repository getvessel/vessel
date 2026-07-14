package models

type NotificationEvent struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	EventType string `json:"eventType"`
	ProjectID string `json:"projectId,omitempty"`
	URL       string `json:"url,omitempty"`
}
