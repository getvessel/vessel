package workspace

import "context"

type Repository interface {
	Create(ctx context.Context, ws *Workspace) error
	Get(ctx context.Context, id string) (*Workspace, error)
	List(ctx context.Context, ownerID string) ([]*Workspace, error)
	Update(ctx context.Context, ws *Workspace) error
	Delete(ctx context.Context, id, ownerID string) error

	CreateTrustedDomain(ctx context.Context, d *TrustedDomain) error
	ListTrustedDomains(ctx context.Context, teamID string) ([]*TrustedDomain, error)
	DeleteTrustedDomain(ctx context.Context, id string) error

	CreateSSHKey(ctx context.Context, key *SSHKey) error
	ListSSHKeys(ctx context.Context, teamID string) ([]*SSHKey, error)
	DeleteSSHKey(ctx context.Context, id string) error

	CreateAuditLog(ctx context.Context, log *AuditLog) error
	ListAuditLogs(ctx context.Context, teamID string, limit int) ([]*AuditLog, error)
}
