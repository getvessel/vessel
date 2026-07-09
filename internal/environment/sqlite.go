package environment

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

func (r *SQLiteRepository) Get(_ context.Context, id string) (*Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, id,
	)
	var env Config
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("environment not found: %s", id)
	}
	if err != nil {
		return nil, err
	}
	env.IsDefault = isDefault == 1
	return &env, nil
}

func (r *SQLiteRepository) ListByProject(_ context.Context, projectID string) ([]Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY is_default DESC, created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}
	defer rows.Close()

	var envs []Config
	for rows.Next() {
		var env Config
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}
		env.IsDefault = isDefault == 1
		envs = append(envs, env)
	}
	return envs, rows.Err()
}

func (r *SQLiteRepository) Create(_ context.Context, env *Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	env.CreatedAt = now
	env.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO environments (id, project_id, name, is_default, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		env.ID, env.ProjectID, env.Name, env.IsDefault, env.CreatedAt, env.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.RUnlock()
	_, err := r.db.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}
