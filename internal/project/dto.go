package project

type CreateProjectRequest struct {
	ID                 string `json:"id"`
	TeamID             string `json:"teamId,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	RepositoryURL      string `json:"repositoryUrl,omitempty"`
	RepositoryURLSnake string `json:"repository_url,omitempty"`
	Branch             string `json:"branch,omitempty"`
	InternalPort       int    `json:"internalPort,omitempty"`
	InternalPortSnake  int    `json:"internal_port,omitempty"`
	Domain             string `json:"domain,omitempty"`
}

type SetEnvVarsRequest map[string]string
