package database

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type SQLiteRepository struct {
	db    *sql.DB
	mu    sync.Mutex
	vault Vault
}

func NewSQLiteRepository(db *sql.DB, vault Vault) *SQLiteRepository {
	return &SQLiteRepository{db: db, vault: vault}
}

func (r *SQLiteRepository) Create(_ context.Context, db *Database) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if db.ID == "" {
		db.ID = uuid.NewString()
	}
	now := time.Now()
	db.CreatedAt = now
	db.UpdatedAt = now

	encryptedPassword, err := r.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT INTO databases (
		id, project_id, environment_id, name, engine, version, port, username, encrypted_password, database_name, volume_path, container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		db.ID, db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.CreatedAt, db.UpdatedAt)
	return err
}

func (r *SQLiteRepository) GetByID(_ context.Context, id string) (*Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var d Database
	var encryptedPassword string

	err := r.db.QueryRow(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at
		FROM databases WHERE id = ?`, id).Scan(
		&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version, &d.Port, &d.Username, &encryptedPassword, &d.DatabaseName, &d.VolumePath, &d.ContainerID, &d.Status, &d.InternalDNS, &d.ExternalDNS, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	plainPassword, err := r.vault.Decrypt(encryptedPassword)
	if err == nil {
		d.Password = plainPassword
	}
	return &d, nil
}

func (r *SQLiteRepository) List(_ context.Context) ([]*Database, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM databases ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Database
	for rows.Next() {
		var d Database
		var encryptedPassword string
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version, &d.Port, &d.Username, &encryptedPassword, &d.DatabaseName, &d.VolumePath, &d.ContainerID, &d.Status, &d.InternalDNS, &d.ExternalDNS, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		if plainPassword, err := r.vault.Decrypt(encryptedPassword); err == nil {
			d.Password = plainPassword
		}
		list = append(list, &d)
	}
	return list, nil
}

func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`DELETE FROM databases WHERE id = ?`, id)
	return err
}

func (r *SQLiteRepository) Update(_ context.Context, db *Database) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	db.UpdatedAt = time.Now()

	encryptedPassword, err := r.vault.Encrypt(db.Password)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`UPDATE databases SET project_id = ?, environment_id = ?, name = ?, engine = ?, version = ?, port = ?, username = ?, encrypted_password = ?, database_name = ?, volume_path = ?, container_id = ?, status = ?, internal_dns = ?, external_dns = ?, updated_at = ? WHERE id = ?`,
		db.ProjectID, db.EnvironmentID, db.Name, db.Engine, db.Version, db.Port, db.Username, encryptedPassword, db.DatabaseName, db.VolumePath, db.ContainerID, db.Status, db.InternalDNS, db.ExternalDNS, db.UpdatedAt, db.ID)
	return err
}
