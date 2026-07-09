package service

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
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Create(_ context.Context, svc *AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if svc.ID == "" {
		svc.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = "building"
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}

	_, err := r.db.Exec(
		`INSERT INTO app_services (id, project_id, environment_id, name, repository_url, branch, internal_port, domain, container_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		svc.ID, svc.ProjectID, svc.EnvironmentID, svc.Name, svc.RepositoryURL, svc.Branch,
		svc.InternalPort, svc.Domain, svc.ContainerID, svc.Status, svc.CreatedAt, svc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create app service: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetByID(_ context.Context, id string) (*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	row := r.db.QueryRow(
		`SELECT id, project_id, environment_id, name, repository_url, branch, internal_port, domain, container_id, status, created_at, updated_at
		FROM app_services WHERE id = ?`, id,
	)
	var svc AppService
	err := row.Scan(
		&svc.ID, &svc.ProjectID, &svc.EnvironmentID, &svc.Name, &svc.RepositoryURL, &svc.Branch,
		&svc.InternalPort, &svc.Domain, &svc.ContainerID, &svc.Status, &svc.CreatedAt, &svc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("app service not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get app service: %w", err)
	}
	return &svc, nil
}

func (r *SQLiteRepository) ListByEnvironment(_ context.Context, environmentID string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(
		`SELECT id, project_id, environment_id, name, repository_url, branch, internal_port, domain, container_id, status, created_at, updated_at
		FROM app_services WHERE environment_id = ? ORDER BY created_at ASC`, environmentID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by environment: %w", err)
	}
	defer rows.Close()

	return scanServices(rows)
}

func (r *SQLiteRepository) ListByProject(_ context.Context, projectID string) ([]*AppService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(
		`SELECT id, project_id, environment_id, name, repository_url, branch, internal_port, domain, container_id, status, created_at, updated_at
		FROM app_services WHERE project_id = ? ORDER BY created_at ASC`, projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list app services by project: %w", err)
	}
	defer rows.Close()

	return scanServices(rows)
}

func (r *SQLiteRepository) Update(_ context.Context, svc *AppService) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	svc.UpdatedAt = time.Now().UTC()

	_, err := r.db.Exec(
		`UPDATE app_services SET
			name = ?, repository_url = ?, branch = ?, internal_port = ?, domain = ?, container_id = ?, status = ?, updated_at = ?
		WHERE id = ?`,
		svc.Name, svc.RepositoryURL, svc.Branch, svc.InternalPort, svc.Domain,
		svc.ContainerID, svc.Status, svc.UpdatedAt, svc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update app service: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`DELETE FROM app_services WHERE id = ?`, id)
	return err
}

func scanServices(rows *sql.Rows) ([]*AppService, error) {
	var list []*AppService
	for rows.Next() {
		var svc AppService
		if err := rows.Scan(
			&svc.ID, &svc.ProjectID, &svc.EnvironmentID, &svc.Name, &svc.RepositoryURL, &svc.Branch,
			&svc.InternalPort, &svc.Domain, &svc.ContainerID, &svc.Status, &svc.CreatedAt, &svc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan app service row: %w", err)
		}
		list = append(list, &svc)
	}
	return list, rows.Err()
}
