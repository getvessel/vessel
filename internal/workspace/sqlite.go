package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Migrate(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			preferred_region TEXT NOT NULL DEFAULT 'local',
			owner_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_trusted_domains (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			domain TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'developer',
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_ssh_keys (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			name TEXT NOT NULL,
			public_key TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspace_audit_logs (
			id TEXT PRIMARY KEY,
			team_id TEXT NOT NULL,
			project_id TEXT DEFAULT '',
			environment_id TEXT DEFAULT '',
			action TEXT NOT NULL,
			actor TEXT NOT NULL,
			created_at TEXT NOT NULL
		);
	`)
	return err
}

func (r *SQLiteRepository) Create(ctx context.Context, ws *Workspace) error {
	if ws.ID == "" {
		ws.ID = uuid.NewString()
	}
	if ws.CreatedAt.IsZero() {
		ws.CreatedAt = time.Now().UTC()
	}
	if ws.UpdatedAt.IsZero() {
		ws.UpdatedAt = time.Now().UTC()
	}
	if ws.PreferredRegion == "" {
		ws.PreferredRegion = "local"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO workspaces (id, name, avatar_url, preferred_region, owner_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ws.ID, ws.Name, ws.AvatarURL, ws.PreferredRegion, ws.OwnerID, ws.CreatedAt.Format(time.RFC3339), ws.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) Get(ctx context.Context, id string) (*Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var ws Workspace
	var createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE id = ?`, id).
		Scan(&ws.ID, &ws.Name, &ws.AvatarURL, &ws.PreferredRegion, &ws.OwnerID, &createdStr, &updatedStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}
	ws.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	ws.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &ws, nil
}

func (r *SQLiteRepository) List(ctx context.Context, ownerID string) ([]*Workspace, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, name, avatar_url, preferred_region, owner_id, created_at, updated_at FROM workspaces WHERE owner_id = ? ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()

	var list []*Workspace
	for rows.Next() {
		var ws Workspace
		var createdStr, updatedStr string
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.AvatarURL, &ws.PreferredRegion, &ws.OwnerID, &createdStr, &updatedStr); err != nil {
			return nil, err
		}
		ws.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		ws.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		list = append(list, &ws)
	}
	return list, nil
}

func (r *SQLiteRepository) Update(ctx context.Context, ws *Workspace) error {
	ws.UpdatedAt = time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `UPDATE workspaces SET name = ?, avatar_url = ?, preferred_region = ?, updated_at = ? WHERE id = ?`,
		ws.Name, ws.AvatarURL, ws.PreferredRegion, ws.UpdatedAt.Format(time.RFC3339), ws.ID)
	if err != nil {
		return fmt.Errorf("update workspace: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int
	_ = r.db.QueryRowContext(ctx, "SELECT count(*) FROM workspaces WHERE owner_id = ?", ownerID).Scan(&count)
	if count <= 1 {
		return errors.New("cannot delete your last workspace. To delete your account, visit Account Settings")
	}

	_, err := r.db.ExecContext(ctx, "DELETE FROM workspaces WHERE id = ? AND owner_id = ?", id, ownerID)
	return err
}

func (r *SQLiteRepository) CreateTrustedDomain(ctx context.Context, d *TrustedDomain) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	if d.Role == "" {
		d.Role = "developer"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_trusted_domains (id, team_id, domain, role, created_at) VALUES (?, ?, ?, ?, ?)`,
		d.ID, d.TeamID, d.Domain, d.Role, d.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create trusted domain: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) ListTrustedDomains(ctx context.Context, teamID string) ([]*TrustedDomain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, team_id, domain, role, created_at FROM workspace_trusted_domains WHERE team_id = ? ORDER BY created_at DESC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list trusted domains: %w", err)
	}
	defer rows.Close()

	var list []*TrustedDomain
	for rows.Next() {
		var d TrustedDomain
		var createdAtStr string
		if err := rows.Scan(&d.ID, &d.TeamID, &d.Domain, &d.Role, &createdAtStr); err != nil {
			return nil, err
		}
		d.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &d)
	}
	return list, nil
}

func (r *SQLiteRepository) DeleteTrustedDomain(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_trusted_domains WHERE id = ?`, id)
	return err
}

func (r *SQLiteRepository) CreateSSHKey(ctx context.Context, key *SSHKey) error {
	if key.ID == "" {
		key.ID = uuid.NewString()
	}
	if key.CreatedAt.IsZero() {
		key.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_ssh_keys (id, team_id, name, public_key, created_at) VALUES (?, ?, ?, ?, ?)`,
		key.ID, key.TeamID, key.Name, key.PublicKey, key.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create ssh key: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) ListSSHKeys(ctx context.Context, teamID string) ([]*SSHKey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, team_id, name, public_key, created_at FROM workspace_ssh_keys WHERE team_id = ? ORDER BY created_at DESC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list ssh keys: %w", err)
	}
	defer rows.Close()

	var list []*SSHKey
	for rows.Next() {
		var k SSHKey
		var createdAtStr string
		if err := rows.Scan(&k.ID, &k.TeamID, &k.Name, &k.PublicKey, &createdAtStr); err != nil {
			return nil, err
		}
		k.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &k)
	}
	return list, nil
}

func (r *SQLiteRepository) DeleteSSHKey(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `DELETE FROM workspace_ssh_keys WHERE id = ?`, id)
	return err
}

func (r *SQLiteRepository) CreateAuditLog(ctx context.Context, log *AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.ExecContext(ctx, `INSERT INTO workspace_audit_logs (id, team_id, project_id, environment_id, action, actor, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.ID, log.TeamID, log.ProjectID, log.EnvironmentID, log.Action, log.Actor, log.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) ListAuditLogs(ctx context.Context, teamID string, limit int) ([]*AuditLog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, team_id, COALESCE(project_id, ''), COALESCE(environment_id, ''), action, actor, created_at FROM workspace_audit_logs WHERE team_id = ? ORDER BY created_at DESC LIMIT ?`, teamID, limit)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	if rows == nil {
		return []*AuditLog{}, nil
	}
	defer rows.Close()

	var list []*AuditLog
	for rows.Next() {
		var log AuditLog
		var createdAtStr string
		if err := rows.Scan(&log.ID, &log.TeamID, &log.ProjectID, &log.EnvironmentID, &log.Action, &log.Actor, &createdAtStr); err != nil {
			return nil, err
		}
		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &log)
	}
	return list, nil
}
