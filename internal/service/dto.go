package service

type CreateAppServiceRequest struct {
	ProjectID     string `json:"projectId"`
	Name          string `json:"name"`
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
	InternalPort  int    `json:"internalPort"`
	Domain        string `json:"domain"`
}

type UpdateAppServiceRequest struct {
	Name          string `json:"name"`
	RepositoryURL string `json:"repositoryUrl"`
	Branch        string `json:"branch"`
	InternalPort  int    `json:"internalPort"`
	Domain        string `json:"domain"`
	ContainerID   string `json:"containerId"`
	Status        string `json:"status"`
}
