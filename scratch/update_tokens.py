import re

def update_repo():
    path = "internal/repositories/project_settings.go"
    with open(path, "r") as f:
        content = f.read()

    # Update interface
    content = content.replace("CreateToken(ctx context.Context, t *models.ProjectToken) (string, error)", "CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error")

    # Update implementation
    old_impl = """func (r *ProjectSettingsSQLiteRepository) CreateToken(ctx context.Context, t *models.ProjectToken) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	t.CreatedAt = time.Now().UTC()
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf("vsl_tok_%s", rawSecret)
	t.TokenPrefix = fullToken[:16]
	scopesStr := strings.Join(t.Scopes, ",")
	ipStr := strings.Join(t.IPAllowlist, ",")
	var expiresAtVal interface{}
	if t.ExpiresAt != nil {
		expiresAtVal = t.ExpiresAt.Format(time.RFC3339)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, scopes, ip_allowlist, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.EnvironmentID, t.Name, t.TokenPrefix, fullToken, scopesStr, ipStr, expiresAtVal, t.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}
	return fullToken, nil
}"""

    new_impl = """func (r *ProjectSettingsSQLiteRepository) CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	scopesStr := strings.Join(t.Scopes, ",")
	ipStr := strings.Join(t.IPAllowlist, ",")
	var expiresAtVal interface{}
	if t.ExpiresAt != nil {
		expiresAtVal = t.ExpiresAt.Format(time.RFC3339)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, scopes, ip_allowlist, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.EnvironmentID, t.Name, t.TokenPrefix, fullToken, scopesStr, ipStr, expiresAtVal, t.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}"""
    content = content.replace(old_impl, new_impl)
    
    with open(path, "w") as f:
        f.write(content)

def update_service():
    path = "internal/services/project_settings_service.go"
    with open(path, "r") as f:
        content = f.read()

    old_func = """func (s *ProjectSettingsService) CreateToken(ctx context.Context, t *models.ProjectToken) (*models.ProjectToken, string, error) {
	if t == nil || t.ProjectID == "" || t.Name == "" {
		return nil, "", errors.New("valid token with projectId and name required")
	}
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	t.CreatedAt = time.Now()
	raw, err := s.repo.CreateToken(ctx, t)
	if err != nil {
		return nil, "", err
	}
	return t, raw, nil
}"""

    new_func = """func (s *ProjectSettingsService) CreateToken(ctx context.Context, t *models.ProjectToken) (*models.ProjectToken, string, error) {
	if t == nil || t.ProjectID == "" || t.Name == "" {
		return nil, "", errors.New("valid token with projectId and name required")
	}
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	t.CreatedAt = time.Now().UTC()
	
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", fmt.Errorf("generate token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf("vsl_tok_%s", rawSecret)
	t.TokenPrefix = fullToken[:16]
	
	err := s.repo.CreateToken(ctx, t, fullToken)
	if err != nil {
		return nil, "", err
	}
	return t, fullToken, nil
}"""
    content = content.replace(old_func, new_func)

    # ensure imports for crypto/rand and encoding/hex
    if '"crypto/rand"' not in content:
        content = content.replace('"context"', '"context"\n\t"crypto/rand"\n\t"encoding/hex"\n\t"fmt"')
    
    with open(path, "w") as f:
        f.write(content)

update_repo()
update_service()
