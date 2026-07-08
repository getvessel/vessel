package types

import "time"

// JobConfig represents a scheduled cron or background worker job associated with a project container.
type JobConfig struct {
	ID         string     `json:"id"`
	ProjectID  string     `json:"projectId"`
	Name       string     `json:"name"`
	Schedule   string     `json:"schedule"` // e.g. "0 * * * *" or "@hourly"
	Command    string     `json:"command"`  // e.g. "php artisan schedule:run"
	Status     string     `json:"status"`   // active, paused, error
	LastRunAt  *time.Time `json:"lastRunAt"`
	LastOutput string     `json:"lastOutput"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}
