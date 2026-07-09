package git

import "context"

type Repository interface {
	SaveProvider(ctx context.Context, gp *GitProviderConfig) error
	GetProvider(ctx context.Context, userID, provider string) (*GitProviderConfig, error)
	GetAnyProviderByType(ctx context.Context, provider string) (*GitProviderConfig, error)
	ListProvidersByUser(ctx context.Context, userID string) ([]*GitProviderConfig, error)
	DeleteProvider(ctx context.Context, userID, provider string) error
}
