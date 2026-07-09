package project

// CreateProjectRequest is the payload for creating a new project with an optional initial application service.
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

// SetEnvVarsRequest is a key-value map used to bulk-update project environment variables.
type SetEnvVarsRequest map[string]string
