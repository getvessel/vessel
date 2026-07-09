package project_settings

import "context"

type Repository interface {
	CreateWebhook(ctx context.Context, w *Webhook) error
	ListWebhooksByProject(ctx context.Context, projectID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, projectID string) error

	CreateToken(ctx context.Context, t *Token) (string, error)
	ListTokensByProject(ctx context.Context, projectID string) ([]*Token, error)
	DeleteToken(ctx context.Context, id, projectID string) error

	AddMember(ctx context.Context, m *ProjectMember) error
	ListMembers(ctx context.Context, projectID string) ([]*ProjectMember, error)
	RemoveMember(ctx context.Context, id, projectID string) error
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID    string
	Email string
}
