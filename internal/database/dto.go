package database

type CreateDatabaseRequest struct {
	ProjectID     string `json:"projectId"`
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
	Engine        string `json:"engine"`
	Version       string `json:"version"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	DatabaseName  string `json:"databaseName"`
	VolumePath    string `json:"volumePath"`
}