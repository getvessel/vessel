package service

import "time"

type AppService struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	RepositoryURL string    `json:"repositoryUrl"`
	Branch        string    `json:"branch"`
	InternalPort  int       `json:"internalPort"`
	Domain        string    `json:"domain"`
	ContainerID   string    `json:"containerId"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
