package canvas

import (
	"time"

	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/storage"
)

type CanvasSummary struct {
	ID                 string              `json:"id"`
	WorkspaceID        string              `json:"workspaceId,omitempty"`
	TeamID             string              `json:"teamId,omitempty"`
	Name               string              `json:"name"`
	Description        string              `json:"description,omitempty"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
	EnvironmentsCount  int                 `json:"environmentsCount"`
	AppsCount          int                 `json:"appsCount"`
	DatabasesCount     int                 `json:"databasesCount"`
	StorageCount       int                 `json:"storageCount"`
	OnlineServices     int                 `json:"onlineServices"`
	TotalServices      int                 `json:"totalServices"`
	ServiceIcons       []string            `json:"serviceIcons"`
	DefaultEnvironment *environment.Config `json:"defaultEnvironment,omitempty"`
}

type EnvironmentCanvas struct {
	Environment *environment.Config   `json:"environment"`
	Apps        []*service.AppService `json:"apps"`
	Databases   []*database.Database  `json:"databases"`
	Storage     []*storage.Storage    `json:"storage"`
}
