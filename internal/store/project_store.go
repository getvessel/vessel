package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
)

// CreateProject inserts a new ProjectConfig record into SQLite and creates its default environment.
func (s *Store) CreateProject(p *types.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO projects (id, workspace_id, team_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.WorkspaceID, p.TeamID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	defaultEnv := &types.EnvironmentConfig{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return s.CreateEnvironment(defaultEnv)
}

// GetProject retrieves a ProjectConfig by its ID from SQLite.
func (s *Store) GetProject(id string) (*types.ProjectConfig, error) {
	row := s.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), team_id, name, description, created_at, updated_at FROM projects WHERE id = ?`, id)

	var p types.ProjectConfig
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// ListProjects retrieves all ProjectConfig records ordered by creation date descending.
func (s *Store) ListProjects() ([]types.ProjectConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(workspace_id, ''), team_id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []types.ProjectConfig
	for rows.Next() {
		var p types.ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// ListProjectsByWorkspace retrieves all ProjectConfig records belonging to a specific workspace ID or team ID.
func (s *Store) ListProjectsByWorkspace(id string) ([]types.ProjectConfig, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(workspace_id, ''), team_id, name, description, created_at, updated_at FROM projects WHERE workspace_id = ? OR team_id = ? ORDER BY created_at DESC`, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []types.ProjectConfig
	for rows.Next() {
		var p types.ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if projects == nil {
		projects = []types.ProjectConfig{}
	}
	return projects, nil
}

// UpdateProject updates an existing ProjectConfig record in SQLite.
func (s *Store) UpdateProject(p *types.ProjectConfig) error {
	p.UpdatedAt = time.Now().UTC()
	_, err := s.db.Exec(`UPDATE projects SET workspace_id = ?, team_id = ?, name = ?, description = ?, updated_at = ? WHERE id = ?`,
		p.WorkspaceID, p.TeamID, p.Name, p.Description, p.UpdatedAt, p.ID,
	)
	return err
}

// DeleteProject deletes a ProjectConfig record from SQLite by ID.
func (s *Store) DeleteProject(id string) error {
	_, err := s.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
