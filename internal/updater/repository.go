package updater

import (
	"context"

	"vessel.dev/vessel/internal/settings"
)

type Repository interface {
	GetServerSettings(ctx context.Context) (*settings.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *settings.ServerSettings) error
}
