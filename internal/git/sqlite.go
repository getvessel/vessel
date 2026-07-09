package git

import (
	"context"
	"database/sql"
	"errors"
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

func (r *SQLiteRepository) SaveProvider(_ context.Context, gp *GitProviderConfig) error {
	if gp.ID == "" {
		gp.ID = uuid.NewString()
	}
	now := time.Now()
	gp.CreatedAt = now
	gp.UpdatedAt = now

	encryptedToken, err := r.vault.Encrypt(gp.AccessToken)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT INTO user_git_providers (id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, provider) DO UPDATE SET encrypted_access_token = excluded.encrypted_access_token, account_name = excluded.account_name, updated_at = excluded.updated_at`,
		gp.ID, gp.UserID, gp.Provider, encryptedToken, gp.AccountName, gp.CreatedAt, gp.UpdatedAt,
	)
	return err
}

func (r *SQLiteRepository) GetProvider(_ context.Context, userID, provider string) (*GitProviderConfig, error) {
	if userID == "" {
		return r.GetAnyProviderByType(context.Background(), provider)
	}

	row := r.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)

	var gp GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	decryptedToken, err := r.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *SQLiteRepository) GetAnyProviderByType(_ context.Context, provider string) (*GitProviderConfig, error) {
	row := r.db.QueryRow(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE provider = ? LIMIT 1`, provider)

	var gp GitProviderConfig
	var encryptedToken string
	err := row.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	decryptedToken, err := r.vault.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}
	gp.AccessToken = decryptedToken
	return &gp, nil
}

func (r *SQLiteRepository) ListProvidersByUser(_ context.Context, userID string) ([]*GitProviderConfig, error) {
	rows, err := r.db.Query(`SELECT id, user_id, provider, encrypted_access_token, account_name, created_at, updated_at
		FROM user_git_providers WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*GitProviderConfig
	for rows.Next() {
		var gp GitProviderConfig
		var encryptedToken string
		if err := rows.Scan(&gp.ID, &gp.UserID, &gp.Provider, &encryptedToken, &gp.AccountName, &gp.CreatedAt, &gp.UpdatedAt); err != nil {
			return nil, err
		}
		decryptedToken, err := r.vault.Decrypt(encryptedToken)
		if err != nil {
			return nil, err
		}
		gp.AccessToken = decryptedToken
		list = append(list, &gp)
	}
	return list, nil
}

func (r *SQLiteRepository) DeleteProvider(_ context.Context, userID, provider string) error {
	_, err := r.db.Exec(`DELETE FROM user_git_providers WHERE user_id = ? AND provider = ?`, userID, provider)
	return err
}
