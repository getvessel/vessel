package service_var

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
	db  *sql.DB
	mu  sync.RWMutex
	svc ServiceRepository
}

func NewSQLiteRepository(db *sql.DB, svc ServiceRepository) *SQLiteRepository {
	return &SQLiteRepository{db: db, svc: svc}
}

func (r *SQLiteRepository) ListByService(_ context.Context, serviceID string) ([]*Variable, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(`SELECT id, service_id, environment_id, key, value, is_secret, created_at, updated_at
		FROM service_vars WHERE service_id = ? ORDER BY key ASC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list service vars: %w", err)
	}
	defer rows.Close()

	var vars []*Variable
	for rows.Next() {
		var v Variable
		var isSecret int
		if err := rows.Scan(&v.ID, &v.ServiceID, &v.EnvironmentID, &v.Key, &v.Value, &isSecret, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan service var: %w", err)
		}
		v.IsSecret = isSecret == 1
		vars = append(vars, &v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vars, nil
}

func (r *SQLiteRepository) Create(_ context.Context, v *Variable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	v.CreatedAt = now
	v.UpdatedAt = now

	isSecret := 0
	if v.IsSecret {
		isSecret = 1
	}

	_, err := r.db.Exec(`INSERT INTO service_vars (id, service_id, environment_id, key, value, is_secret, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(service_id, key) DO UPDATE SET
		value = excluded.value,
		is_secret = excluded.is_secret,
		updated_at = excluded.updated_at`,
		v.ID, v.ServiceID, v.EnvironmentID, v.Key, v.Value, isSecret, v.CreatedAt, v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create service var: %w", err)
	}

	_, _ = r.db.Exec(`UPDATE app_services SET env_vars_count = (SELECT COUNT(*) FROM service_vars WHERE service_id = ?) WHERE id = ?`, v.ServiceID, v.ServiceID)
	return nil
}

func (r *SQLiteRepository) Update(_ context.Context, v *Variable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	v.UpdatedAt = time.Now().UTC()

	isSecret := 0
	if v.IsSecret {
		isSecret = 1
	}

	res, err := r.db.Exec(`UPDATE service_vars SET key = ?, value = ?, is_secret = ?, updated_at = ?
		WHERE id = ? AND service_id = ?`,
		v.Key, v.Value, isSecret, v.UpdatedAt, v.ID, v.ServiceID)
	if err != nil {
		return fmt.Errorf("update service var: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("service variable not found")
	}

	_, _ = r.db.Exec(`UPDATE app_services SET env_vars_count = (SELECT COUNT(*) FROM service_vars WHERE service_id = ?) WHERE id = ?`, v.ServiceID, v.ServiceID)
	return nil
}

func (r *SQLiteRepository) Delete(_ context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.db.Exec(`DELETE FROM service_vars WHERE id = ? AND service_id = ?`, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete service var: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("service variable not found")
	}

	_, _ = r.db.Exec(`UPDATE app_services SET env_vars_count = (SELECT COUNT(*) FROM service_vars WHERE service_id = ?) WHERE id = ?`, serviceID, serviceID)
	return nil
}

func (r *SQLiteRepository) BulkSet(_ context.Context, serviceID, environmentID string, vars []*Variable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM service_vars WHERE service_id = ?`, serviceID); err != nil {
		return fmt.Errorf("bulk delete: %w", err)
	}

	now := time.Now().UTC()
	for _, v := range vars {
		if v.ID == "" {
			v.ID = uuid.NewString()
		}
		v.ServiceID = serviceID
		v.EnvironmentID = environmentID
		v.CreatedAt = now
		v.UpdatedAt = now

		isSecret := 0
		if v.IsSecret {
			isSecret = 1
		}

		_, err := tx.Exec(`INSERT INTO service_vars (id, service_id, environment_id, key, value, is_secret, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			v.ID, v.ServiceID, v.EnvironmentID, v.Key, v.Value, isSecret, v.CreatedAt, v.UpdatedAt)
		if err != nil {
			return fmt.Errorf("bulk insert: %w", err)
		}
	}

	if _, err := tx.Exec(`UPDATE app_services SET env_vars_count = ? WHERE id = ?`, len(vars), serviceID); err != nil {
		return fmt.Errorf("update env_vars_count: %w", err)
	}

	return tx.Commit()
}
