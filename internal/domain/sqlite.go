package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) ListByProject(_ context.Context, projectID string) ([]Config, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains WHERE project_id = ? ORDER BY domain_name ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []Config
	for rows.Next() {
		var d Config
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *SQLiteRepository) ListAll(ctx context.Context) ([]Config, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains ORDER BY domain_name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []Config
	for rows.Next() {
		var d Config
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *SQLiteRepository) Create(_ context.Context, d *Config) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO domains (id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, d.DomainName, d.RedirectTo, d.SSLCertStatus, d.PathPrefix, d.CreatedAt, d.UpdatedAt,
	)
	return err
}

func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}
