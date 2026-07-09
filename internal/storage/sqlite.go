package storage

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

const listStorageQuery = `SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, COALESCE(volume_path, ''), COALESCE(container_id, ''), COALESCE(status, 'stopped'), COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage`

func scanStorage(scanner interface {
	Scan(dest ...any) error
}, s *Storage, encryptedSecretKey *string) error {
	return scanner.Scan(
		&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Name, &s.Type,
		&s.APIPort, &s.ConsolePort, &s.AccessKey, encryptedSecretKey,
		&s.BucketName, &s.VolumePath, &s.ContainerID, &s.Status,
		&s.InternalDNS, &s.ExternalDNS, &s.CreatedAt, &s.UpdatedAt,
	)
}

func (r *SQLiteRepository) decryptSecretKey(encrypted string, s *Storage) {
	if plain, err := r.vault.Decrypt(encrypted); err == nil {
		s.SecretKey = plain
	}
}

func (r *SQLiteRepository) Create(_ context.Context, s *Storage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	encryptedSecretKey, err := r.vault.Encrypt(s.SecretKey)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT INTO storage (
		id, project_id, environment_id, name, type, api_port, console_port,
		access_key, encrypted_secret_key, bucket_name, volume_path,
		container_id, status, internal_dns, external_dns, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		s.ID, s.ProjectID, s.EnvironmentID, s.Name, s.Type,
		s.APIPort, s.ConsolePort, s.AccessKey, encryptedSecretKey,
		s.BucketName, s.VolumePath, s.ContainerID, s.Status,
		s.InternalDNS, s.ExternalDNS, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *SQLiteRepository) GetByID(_ context.Context, id string) (*Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := r.db.QueryRow(listStorageQuery+` WHERE id = ?`, id)

	var s Storage
	var encryptedSecretKey string
	if err := scanStorage(row, &s, &encryptedSecretKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	r.decryptSecretKey(encryptedSecretKey, &s)
	return &s, nil
}

func (r *SQLiteRepository) List(_ context.Context) ([]*Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(listStorageQuery + ` ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Storage
	for rows.Next() {
		var s Storage
		var encryptedSecretKey string
		if err := scanStorage(rows, &s, &encryptedSecretKey); err != nil {
			return nil, err
		}
		r.decryptSecretKey(encryptedSecretKey, &s)
		list = append(list, &s)
	}
	return list, nil
}

func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`DELETE FROM storage WHERE id = ?`, id)
	return err
}
