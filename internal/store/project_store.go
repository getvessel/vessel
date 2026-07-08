package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/solomonolatunji/vessel/internal/types"
)

// CreateProject inserts a new ProjectConfig record into SQLite.
func (s *Store) CreateProject(p *types.ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := s.db.Exec(`INSERT INTO projects (
		id, name, repository_url, branch, build_command, start_command, dockerfile_path,
		internal_port, domain, auto_deploy_webhook, cpu_request, memory_limit_mb, health_check_path,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.RepositoryURL, p.Branch, p.BuildCommand, p.StartCommand, p.DockerfilePath,
		p.InternalPort, p.Domain, p.AutoDeployWebhook, p.CPURequest, p.MemoryLimitMB, p.HealthCheckPath,
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

// GetProject retrieves a ProjectConfig record by its ID.
func (s *Store) GetProject(id string) (*types.ProjectConfig, error) {
	row := s.db.QueryRow(`SELECT id, name, repository_url, branch, build_command, start_command, dockerfile_path,
		internal_port, domain, auto_deploy_webhook, cpu_request, memory_limit_mb, health_check_path, created_at, updated_at
		FROM projects WHERE id = ?`, id)

	var p types.ProjectConfig
	err := row.Scan(&p.ID, &p.Name, &p.RepositoryURL, &p.Branch, &p.BuildCommand, &p.StartCommand, &p.DockerfilePath,
		&p.InternalPort, &p.Domain, &p.AutoDeployWebhook, &p.CPURequest, &p.MemoryLimitMB, &p.HealthCheckPath,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListProjects retrieves all ProjectConfig records ordered by creation date descending.
func (s *Store) ListProjects() ([]types.ProjectConfig, error) {
	rows, err := s.db.Query(`SELECT id, name, repository_url, branch, build_command, start_command, dockerfile_path,
		internal_port, domain, auto_deploy_webhook, cpu_request, memory_limit_mb, health_check_path, created_at, updated_at
		FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []types.ProjectConfig
	for rows.Next() {
		var p types.ProjectConfig
		if err := rows.Scan(&p.ID, &p.Name, &p.RepositoryURL, &p.Branch, &p.BuildCommand, &p.StartCommand, &p.DockerfilePath,
			&p.InternalPort, &p.Domain, &p.AutoDeployWebhook, &p.CPURequest, &p.MemoryLimitMB, &p.HealthCheckPath,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// DeleteProject deletes a ProjectConfig record from SQLite by ID.
func (s *Store) DeleteProject(id string) error {
	_, err := s.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
