package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateWorkspaceTrustedDomain inserts a new trusted domain entry for a workspace team.
func (s *Store) CreateWorkspaceTrustedDomain(item *types.WorkspaceTrustedDomain) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	if item.Role == "" {
		item.Role = "developer"
	}

	query := `INSERT INTO workspace_trusted_domains (id, team_id, domain, role, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, item.ID, item.TeamID, item.Domain, item.Role, item.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace trusted domain: %w", err)
	}
	return nil
}

// ListWorkspaceTrustedDomains returns all trusted domains for a given team workspace.
func (s *Store) ListWorkspaceTrustedDomains(teamID string) ([]*types.WorkspaceTrustedDomain, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, team_id, domain, role, created_at FROM workspace_trusted_domains WHERE team_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace trusted domains: %w", err)
	}
	defer rows.Close()

	var list []*types.WorkspaceTrustedDomain
	for rows.Next() {
		var item types.WorkspaceTrustedDomain
		var createdAtStr string
		if err := rows.Scan(&item.ID, &item.TeamID, &item.Domain, &item.Role, &createdAtStr); err != nil {
			return nil, err
		}
		item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &item)
	}
	return list, nil
}

// DeleteWorkspaceTrustedDomain removes a trusted domain by ID.
func (s *Store) DeleteWorkspaceTrustedDomain(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM workspace_trusted_domains WHERE id = ?`, id)
	return err
}

// CreateWorkspaceSSHKey inserts a new SSH public key entry for a workspace team.
func (s *Store) CreateWorkspaceSSHKey(item *types.WorkspaceSSHKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO workspace_ssh_keys (id, team_id, name, public_key, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, item.ID, item.TeamID, item.Name, item.PublicKey, item.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace ssh key: %w", err)
	}
	return nil
}

// ListWorkspaceSSHKeys returns all SSH public keys for a given team workspace.
func (s *Store) ListWorkspaceSSHKeys(teamID string) ([]*types.WorkspaceSSHKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `SELECT id, team_id, name, public_key, created_at FROM workspace_ssh_keys WHERE team_id = ? ORDER BY created_at DESC`
	rows, err := s.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace ssh keys: %w", err)
	}
	defer rows.Close()

	var list []*types.WorkspaceSSHKey
	for rows.Next() {
		var item types.WorkspaceSSHKey
		var createdAtStr string
		if err := rows.Scan(&item.ID, &item.TeamID, &item.Name, &item.PublicKey, &createdAtStr); err != nil {
			return nil, err
		}
		item.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &item)
	}
	return list, nil
}

// DeleteWorkspaceSSHKey removes an SSH key by ID.
func (s *Store) DeleteWorkspaceSSHKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM workspace_ssh_keys WHERE id = ?`, id)
	return err
}

// CreateWorkspaceAuditLog inserts an audit log event into the workspace history.
func (s *Store) CreateWorkspaceAuditLog(log *types.WorkspaceAuditLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO workspace_audit_logs (id, team_id, project_id, environment_id, action, actor, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, log.ID, log.TeamID, log.ProjectID, log.EnvironmentID, log.Action, log.Actor, log.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to create workspace audit log: %w", err)
	}
	return nil
}

// ListWorkspaceAuditLogs returns the audit history for a given team workspace.
func (s *Store) ListWorkspaceAuditLogs(teamID string, limit int) ([]*types.WorkspaceAuditLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, team_id, COALESCE(project_id, ''), COALESCE(environment_id, ''), action, actor, created_at
		FROM workspace_audit_logs WHERE team_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := s.db.Query(query, teamID, limit)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query workspace audit logs: %w", err)
	}
	if rows == nil {
		return []*types.WorkspaceAuditLog{}, nil
	}
	defer rows.Close()

	var list []*types.WorkspaceAuditLog
	for rows.Next() {
		var log types.WorkspaceAuditLog
		var createdAtStr string
		if err := rows.Scan(&log.ID, &log.TeamID, &log.ProjectID, &log.EnvironmentID, &log.Action, &log.Actor, &createdAtStr); err != nil {
			return nil, err
		}
		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		list = append(list, &log)
	}
	return list, nil
}
