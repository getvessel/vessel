package repositories

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
)

type ProjectSettingsRepository interface {
	CreateWebhook(ctx context.Context, w *models.Webhook) error
	ListWebhooksByProject(ctx context.Context, projectID string) ([]*models.Webhook, error)
	DeleteWebhook(ctx context.Context, id, projectID string) error

	CreateToken(ctx context.Context, t *models.ProjectToken) (string, error)
	ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error)
	DeleteToken(ctx context.Context, id, projectID string) error

	AddMember(ctx context.Context, m *models.ProjectMember) error
	ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error)
	RemoveMember(ctx context.Context, id, projectID string) error
}

type ProjectSettingsSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewProjectSettingsSQLiteRepository(db *sql.DB) *ProjectSettingsSQLiteRepository {
	return &ProjectSettingsSQLiteRepository{db: db}
}

func (r *ProjectSettingsSQLiteRepository) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now

	eventTypesStr := strings.Join(w.EventTypes, ",")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_webhooks (id, project_id, url, event_types, include_pr_environments, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		w.ID, w.ProjectID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) ListWebhooksByProject(ctx context.Context, projectID string) ([]*models.Webhook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, url, event_types, include_pr_environments, created_at, updated_at
		 FROM project_webhooks WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	defer rows.Close()

	var out []*models.Webhook
	for rows.Next() {
		var w models.Webhook
		var eventsStr string
		var includePr int
		if err := rows.Scan(&w.ID, &w.ProjectID, &w.URL, &eventsStr, &includePr, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan webhook: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		w.IncludePREnvironments = includePr == 1
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsSQLiteRepository) DeleteWebhook(ctx context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `DELETE FROM project_webhooks WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("webhook not found")
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) CreateToken(ctx context.Context, t *models.ProjectToken) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	t.CreatedAt = time.Now().UTC()

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf("vsl_tok_%s", rawSecret)
	t.TokenPrefix = fullToken[:16]

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.EnvironmentID, t.Name, t.TokenPrefix, fullToken, t.CreatedAt)
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}
	return fullToken, nil
}

func (r *ProjectSettingsSQLiteRepository) ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, environment_id, name, token_prefix, created_at
		 FROM project_tokens WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()

	var out []*models.ProjectToken
	for rows.Next() {
		var t models.ProjectToken
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsSQLiteRepository) DeleteToken(ctx context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `DELETE FROM project_tokens WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("token not found")
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) AddMember(ctx context.Context, m *models.ProjectMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	m.InvitedAt = now
	if m.Status == "" {
		m.Status = "pending"
	}
	if m.Permission == "" {
		m.Permission = "Can Edit"
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_members (id, project_id, user_id, email, permission, status, invited_at, accepted_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, email) DO UPDATE SET
		 permission = excluded.permission,
		 status = excluded.status`,
		m.ID, m.ProjectID, m.UserID, m.Email, m.Permission, m.Status, m.InvitedAt, m.AcceptedAt)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r *ProjectSettingsSQLiteRepository) ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, user_id, email, permission, status, invited_at, accepted_at
		 FROM project_members WHERE project_id = ? ORDER BY invited_at ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	var out []*models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		var acceptedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Email, &m.Permission, &m.Status, &m.InvitedAt, &acceptedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		if acceptedAt.Valid {
			m.AcceptedAt = acceptedAt.Time
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsSQLiteRepository) RemoveMember(ctx context.Context, id, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.ExecContext(ctx, `DELETE FROM project_members WHERE id = ? AND project_id = ?`, id, projectID)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("member not found")
	}
	return nil
}
