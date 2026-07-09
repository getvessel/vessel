package store

import (
	"database/sql"
	"fmt"
	"time"

	"vessel.dev/vessel/internal/types"
)

func (s *Store) ListOAuthProviders() ([]types.OAuthProvider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, provider_name, enabled, COALESCE(client_id, ''), COALESCE(client_secret, ''), COALESCE(redirect_uri, ''), COALESCE(base_url, ''), COALESCE(tenant, ''), created_at, updated_at FROM oauth_providers ORDER BY provider_name ASC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list oauth providers: %w", err)
	}
	defer rows.Close()

	var providers []types.OAuthProvider
	for rows.Next() {
		var p types.OAuthProvider
		if err := rows.Scan(&p.ID, &p.ProviderName, &p.Enabled, &p.ClientID, &p.ClientSecret, &p.RedirectURI, &p.BaseURL, &p.Tenant, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed scanning oauth provider row: %w", err)
		}
		providers = append(providers, p)
	}

	return providers, nil
}

func (s *Store) GetOAuthProvider(idOrName string) (*types.OAuthProvider, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `SELECT id, provider_name, enabled, COALESCE(client_id, ''), COALESCE(client_secret, ''), COALESCE(redirect_uri, ''), COALESCE(base_url, ''), COALESCE(tenant, ''), created_at, updated_at FROM oauth_providers WHERE id = ? OR provider_name = ?`
	row := s.db.QueryRow(query, idOrName, idOrName)

	var p types.OAuthProvider
	if err := row.Scan(&p.ID, &p.ProviderName, &p.Enabled, &p.ClientID, &p.ClientSecret, &p.RedirectURI, &p.BaseURL, &p.Tenant, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("oauth provider not found: %s", idOrName)
		}
		return nil, fmt.Errorf("failed to scan oauth provider: %w", err)
	}

	return &p, nil
}

func (s *Store) SaveOAuthProvider(p *types.OAuthProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	query := `INSERT INTO oauth_providers (
		id, provider_name, enabled, client_id, client_secret, redirect_uri, base_url, tenant, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		provider_name = excluded.provider_name,
		enabled = excluded.enabled,
		client_id = excluded.client_id,
		client_secret = excluded.client_secret,
		redirect_uri = excluded.redirect_uri,
		base_url = excluded.base_url,
		tenant = excluded.tenant,
		updated_at = excluded.updated_at`

	_, err := s.db.Exec(query, p.ID, p.ProviderName, p.Enabled, p.ClientID, p.ClientSecret, p.RedirectURI, p.BaseURL, p.Tenant, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save oauth provider: %w", err)
	}

	return nil
}

func (s *Store) GetUserTOTPSecret(userID string) (string, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var secret string
	var recovery string
	err := s.db.QueryRow(`SELECT COALESCE(totp_secret, ''), COALESCE(recovery_codes, '') FROM users WHERE id = ?`, userID).Scan(&secret, &recovery)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get totp secret: %w", err)
	}

	var codes []string
	if recovery != "" {
		// Split comma-separated recovery codes
		var current string
		for _, r := range recovery {
			if r == ',' {
				if current != "" {
					codes = append(codes, current)
					current = ""
				}
			} else {
				current += string(r)
			}
		}
		if current != "" {
			codes = append(codes, current)
		}
	}

	return secret, codes, nil
}

func (s *Store) UpdateUserTOTP(userID string, enabled bool, secret string, recoveryCodes []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	recoveryStr := ""
	for i, c := range recoveryCodes {
		if i > 0 {
			recoveryStr += ","
		}
		recoveryStr += c
	}

	query := `UPDATE users SET totp_enabled = ?, totp_secret = ?, recovery_codes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := s.db.Exec(query, enabled, secret, recoveryStr, userID)
	if err != nil {
		return fmt.Errorf("failed to update user totp status: %w", err)
	}

	return nil
}
