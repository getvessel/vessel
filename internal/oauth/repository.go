package oauth

import "context"

type Repository interface {
	ListProviders(ctx context.Context) ([]Provider, error)
	GetProvider(ctx context.Context, idOrName string) (*Provider, error)
	SaveProvider(ctx context.Context, p *Provider) error

	GetUserTOTPSecret(ctx context.Context, userID string) (secret string, recoveryCodes []string, err error)
	UpdateUserTOTP(ctx context.Context, userID string, enabled bool, secret string, recoveryCodes []string) error
}
