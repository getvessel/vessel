package models

type RailwayImportRequest struct {
	Token              string `json:"token"`
	ProjectID          string `json:"projectId"`
	ExcludeRailwayVars bool   `json:"excludeRailwayVars"`
	RecreateDatabases  bool   `json:"recreateDatabases"`
	ImportData         bool   `json:"importData"`
}

type RailwayProjectsResponse struct {
	Data RailwayProjectsData `json:"data"`
}

type RailwayProjectsData struct {
	Projects RailwayProjectConnection `json:"projects"`
}

type RailwayProjectConnection struct {
	Edges []RailwayProjectEdge `json:"edges"`
}

type RailwayProjectEdge struct {
	Node RailwayProject `json:"node"`
}

type RailwayProject struct {
	ID           string                       `json:"id"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Environments RailwayEnvironmentConnection `json:"environments"`
	Services     RailwayServiceConnection     `json:"services"`
}

type RailwayEnvironmentConnection struct {
	Edges []RailwayEnvironmentEdge `json:"edges"`
}

type RailwayEnvironmentEdge struct {
	Node RailwayEnvironment `json:"node"`
}

type RailwayEnvironment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RailwayServiceConnection struct {
	Edges []RailwayServiceEdge `json:"edges"`
}

type RailwayServiceEdge struct {
	Node RailwayService `json:"node"`
}

type RailwayService struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	Source RailwayServiceSource `json:"source"`
}

type RailwayServiceSource struct {
	Image string `json:"image"`
	Repo  string `json:"repo"`
}

type RailwayProjectDetailsResponse struct {
	Data RailwayProjectDetailsData `json:"data"`
}

type RailwayProjectDetailsData struct {
	Project RailwayProject `json:"project"`
}
