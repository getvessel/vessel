package updater

import (
	"context"

	"vessel.dev/vessel/internal/settings"
)

// Repository defines the data operations required by UpdaterService ("Accept interfaces, return structs").
type Repository interface {
	GetServerSettings(ctx context.Context) (*settings.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *settings.ServerSettings) error
}
