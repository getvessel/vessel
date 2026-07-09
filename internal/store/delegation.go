package store

import (
	"context"

	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/user"
)

// ── Settings ─────────────────────────────────────────────────────────────────

// GetServerSettings delegates to the modular settings repository.
func (s *Store) GetServerSettings() (*settings.ServerSettings, error) {
	return settings.NewSQLiteRepository(s.db).GetServerSettings(context.Background())
}

// UpdateServerSettings delegates to the modular settings repository.
func (s *Store) UpdateServerSettings(cfg *settings.ServerSettings) error {
	return settings.NewSQLiteRepository(s.db).UpdateServerSettings(context.Background(), cfg)
}

// ── User ─────────────────────────────────────────────────────────────────────

// CreateUser delegates to the modular user repository.
func (s *Store) CreateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).CreateUser(context.Background(), u)
}

// GetUserByEmail delegates to the modular user repository.
func (s *Store) GetUserByEmail(email string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByEmail(context.Background(), email)
}

// GetUserByID delegates to the modular user repository.
func (s *Store) GetUserByID(id string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByID(context.Background(), id)
}

// ListUsers delegates to the modular user repository.
func (s *Store) ListUsers() ([]user.User, error) {
	return user.NewSQLiteRepository(s.db).ListUsers(context.Background())
}

// UpdateUser delegates to the modular user repository.
func (s *Store) UpdateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).UpdateUser(context.Background(), u)
}

// CreatePersonalAccessToken delegates to the modular user repository.
func (s *Store) CreatePersonalAccessToken(pat *user.PersonalAccessToken) error {
	return user.NewSQLiteRepository(s.db).CreatePAT(context.Background(), pat)
}

// ListPersonalAccessTokens delegates to the modular user repository.
func (s *Store) ListPersonalAccessTokens(userID string) ([]*user.PersonalAccessToken, error) {
	return user.NewSQLiteRepository(s.db).ListPATs(context.Background(), userID)
}

// DeletePersonalAccessToken delegates to the modular user repository.
func (s *Store) DeletePersonalAccessToken(id, userID string) error {
	return user.NewSQLiteRepository(s.db).DeletePAT(context.Background(), id, userID)
}

// ── Projects ─────────────────────────────────────────────────────────────────

func (s *Store) projectRepo() project.Repository {
	return project.NewSQLiteRepository(s.db, s.vault)
}

// ListProjects delegates to the modular project repository.
func (s *Store) ListProjects() ([]types.ProjectConfig, error) {
	projects, err := s.projectRepo().ListProjects(context.Background())
	if err != nil {
		return nil, err
	}
	result := make([]types.ProjectConfig, len(projects))
	for i, p := range projects {
		result[i] = toTypesProjectConfig(p)
	}
	return result, nil
}

// ListProjectsByWorkspace delegates to the modular project repository.
func (s *Store) ListProjectsByWorkspace(workspaceID string) ([]types.ProjectConfig, error) {
	all, err := s.projectRepo().ListProjects(context.Background())
	if err != nil {
		return nil, err
	}
	var result []types.ProjectConfig
	for _, p := range all {
		if p.WorkspaceID == workspaceID {
			result = append(result, toTypesProjectConfig(p))
		}
	}
	return result, nil
}

// GetProject delegates to the modular project repository.
func (s *Store) GetProject(id string) (*types.ProjectConfig, error) {
	p, err := s.projectRepo().GetProject(context.Background(), id)
	if err != nil || p == nil {
		return nil, err
	}
	cfg := toTypesProjectConfig(*p)
	return &cfg, nil
}

// CreateProject delegates to the modular project repository.
func (s *Store) CreateProject(p *types.ProjectConfig) error {
	cfg := project.ProjectConfig{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		TeamID:      p.TeamID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	return s.projectRepo().CreateProject(context.Background(), &cfg)
}

// DeleteProject delegates to the modular project repository.
func (s *Store) DeleteProject(id string) error {
	return s.projectRepo().DeleteProject(context.Background(), id)
}

func toTypesProjectConfig(p project.ProjectConfig) types.ProjectConfig {
	return types.ProjectConfig{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		TeamID:      p.TeamID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ── Environments ─────────────────────────────────────────────────────────────

// ListEnvironments delegates to the modular project repository.
func (s *Store) ListEnvironments(projectID string) ([]types.EnvironmentConfig, error) {
	envs, err := s.projectRepo().ListEnvironments(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	result := make([]types.EnvironmentConfig, len(envs))
	for i, e := range envs {
		result[i] = types.EnvironmentConfig{
			ID:        e.ID,
			ProjectID: e.ProjectID,
			Name:      e.Name,
			IsDefault: e.IsDefault,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}
	return result, nil
}

// CreateEnvironment delegates to the modular project repository.
func (s *Store) CreateEnvironment(env *types.EnvironmentConfig) error {
	cfg := project.EnvironmentConfig{
		ID:        env.ID,
		ProjectID: env.ProjectID,
		Name:      env.Name,
		IsDefault: env.IsDefault,
		CreatedAt: env.CreatedAt,
		UpdatedAt: env.UpdatedAt,
	}
	return s.projectRepo().CreateEnvironment(context.Background(), &cfg)
}

// DeleteEnvironment delegates to the modular project repository.
func (s *Store) DeleteEnvironment(id string) error {
	return s.projectRepo().DeleteEnvironment(context.Background(), id)
}

// ── Domains ──────────────────────────────────────────────────────────────────

// ListDomains delegates to the modular project repository.
func (s *Store) ListDomains(projectID string) ([]types.DomainConfig, error) {
	domains, err := s.projectRepo().ListDomains(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	result := make([]types.DomainConfig, len(domains))
	for i, d := range domains {
		result[i] = toTypesDomainConfig(d)
	}
	return result, nil
}

// ListAllDomains returns every custom domain across all projects.
func (s *Store) ListAllDomains() ([]types.DomainConfig, error) {
	// Query all projects and aggregate their domains; projects without domains are skipped.
	projects, err := s.projectRepo().ListProjects(context.Background())
	if err != nil {
		return nil, err
	}
	var result []types.DomainConfig
	for _, p := range projects {
		domains, err := s.projectRepo().ListDomains(context.Background(), p.ID)
		if err != nil {
			return nil, err
		}
		for _, d := range domains {
			result = append(result, toTypesDomainConfig(d))
		}
	}
	return result, nil
}

// AddDomain delegates to the modular project repository.
func (s *Store) AddDomain(d *types.DomainConfig) error {
	cfg := project.DomainConfig{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		DomainName:    d.DomainName,
		RedirectTo:    d.RedirectTo,
		SSLCertStatus: d.SSLCertStatus,
		PathPrefix:    d.PathPrefix,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
	return s.projectRepo().AddDomain(context.Background(), &cfg)
}

// DeleteDomain delegates to the modular project repository.
func (s *Store) DeleteDomain(id string) error {
	return s.projectRepo().DeleteDomain(context.Background(), id)
}

func toTypesDomainConfig(d project.DomainConfig) types.DomainConfig {
	return types.DomainConfig{
		ID:            d.ID,
		ProjectID:     d.ProjectID,
		DomainName:    d.DomainName,
		RedirectTo:    d.RedirectTo,
		SSLCertStatus: d.SSLCertStatus,
		PathPrefix:    d.PathPrefix,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// ── Env Vars ─────────────────────────────────────────────────────────────────

// GetEnvVars delegates to the modular project repository.
func (s *Store) GetEnvVars(projectID string) (map[string]string, error) {
	return s.projectRepo().GetEnvVars(context.Background(), projectID)
}

// SetEnvVar delegates to the modular project repository.
func (s *Store) SetEnvVar(projectID, key, value string) error {
	return s.projectRepo().SetEnvVar(context.Background(), projectID, key, value)
}

// ── Canvas ───────────────────────────────────────────────────────────────────

// ListProjectCanvasSummaries delegates to the modular project repository.
func (s *Store) ListProjectCanvasSummaries() ([]types.ProjectCanvasSummary, error) {
	summaries, err := s.projectRepo().ListProjectCanvasSummaries(context.Background())
	if err != nil {
		return nil, err
	}
	result := make([]types.ProjectCanvasSummary, len(summaries))
	for i, summary := range summaries {
		result[i] = toTypesProjectCanvasSummary(summary)
	}
	return result, nil
}

// GetProjectCanvasSummary delegates to the modular project repository.
func (s *Store) GetProjectCanvasSummary(id string) (*types.ProjectCanvasSummary, error) {
	summary, err := s.projectRepo().GetProjectCanvasSummary(context.Background(), id)
	if err != nil || summary == nil {
		return nil, err
	}
	result := toTypesProjectCanvasSummary(*summary)
	return &result, nil
}

// GetEnvironmentCanvas delegates to the modular project repository.
func (s *Store) GetEnvironmentCanvas(id string) (*types.EnvironmentCanvas, error) {
	canvas, err := s.projectRepo().GetEnvironmentCanvas(context.Background(), id)
	if err != nil || canvas == nil {
		return nil, err
	}
	env := types.EnvironmentConfig{
		ID:        canvas.Environment.ID,
		ProjectID: canvas.Environment.ProjectID,
		Name:      canvas.Environment.Name,
		IsDefault: canvas.Environment.IsDefault,
		CreatedAt: canvas.Environment.CreatedAt,
		UpdatedAt: canvas.Environment.UpdatedAt,
	}
	result := &types.EnvironmentCanvas{
		Environment: &env,
		Apps:        canvas.Apps,
		Databases:   canvas.Databases,
		Storage:     canvas.Storage,
	}
	return result, nil
}

func toTypesProjectCanvasSummary(summary project.ProjectCanvasSummary) types.ProjectCanvasSummary {
	return types.ProjectCanvasSummary{
		ProjectConfig:      toTypesProjectConfig(summary.ProjectConfig),
		EnvironmentsCount:  summary.EnvironmentsCount,
		AppsCount:          summary.AppsCount,
		DatabasesCount:     summary.DatabasesCount,
		StorageCount:       summary.StorageCount,
		OnlineServices:     summary.OnlineServices,
		TotalServices:      summary.TotalServices,
		ServiceIcons:       summary.ServiceIcons,
		DefaultEnvironment: (*types.EnvironmentConfig)(summary.DefaultEnvironment),
	}
}
