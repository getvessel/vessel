package service_var

import "time"

type Variable struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"serviceId"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	IsSecret      bool      `json:"isSecret"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
