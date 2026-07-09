package env

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type SQLiteRepository struct {
	db    *sql.DB
	vault Vault
}

func NewSQLiteRepository(db *sql.DB, vault Vault) *SQLiteRepository {
	return &SQLiteRepository{db: db, vault: vault}
}

func (r *SQLiteRepository) GetVars(_ context.Context, projectID string) (map[string]string, error) {
	rows, err := r.db.Query(`SELECT key, encrypted_value FROM env_vars WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	envs := make(map[string]string)
	for rows.Next() {
		var key, encrypted string
		if err := rows.Scan(&key, &encrypted); err != nil {
			return nil, err
		}
		plaintext, err := r.vault.Decrypt(encrypted)
		if err != nil {
			continue
		}
		envs[key] = plaintext
	}
	return envs, rows.Err()
}

func (r *SQLiteRepository) SetVar(_ context.Context, projectID, key, plaintextValue string) error {
	encrypted, err := r.vault.Encrypt(plaintextValue)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO env_vars (id, project_id, key, encrypted_value, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(project_id, key) DO UPDATE SET encrypted_value = excluded.encrypted_value, updated_at = excluded.updated_at`,
		uuid.NewString(), projectID, key, encrypted, now, now,
	)
	return err
}
