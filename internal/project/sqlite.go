package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
)

// Vault is the minimal interface needed to encrypt/decrypt env-var values.
type Vault interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// appServiceRow is a minimal read-only projection used by canvas queries.
type appServiceRow struct {
	ProjectID     string
	EnvironmentID string
	Status        string
}

// SQLiteRepository implements Repository against a SQLite database.
type SQLiteRepository struct {
	mu    sync.RWMutex
	db    *sql.DB
	vault Vault
}

// NewSQLiteRepository constructs a SQLiteRepository backed by the given db and vault.
func NewSQLiteRepository(db *sql.DB, vault Vault) *SQLiteRepository {
	return &SQLiteRepository{db: db, vault: vault}
}

// ── Projects ─────────────────────────────────────────────────────────────────

// ListProjects retrieves all ProjectConfig records ordered by creation date descending.
func (r *SQLiteRepository) ListProjects(_ context.Context) ([]ProjectConfig, error) {
	rows, err := r.db.Query(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectConfig
	for rows.Next() {
		var p ProjectConfig
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// GetProject retrieves a single ProjectConfig by its ID.
func (r *SQLiteRepository) GetProject(_ context.Context, id string) (*ProjectConfig, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p ProjectConfig
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateProject inserts a new project and creates its default production environment.
func (r *SQLiteRepository) CreateProject(ctx context.Context, p *ProjectConfig) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := r.db.Exec(
		`INSERT INTO projects (id, workspace_id, team_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.WorkspaceID, p.TeamID, p.Name, p.Description, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	defaultEnv := &EnvironmentConfig{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.CreateEnvironment(ctx, defaultEnv)
}

// DeleteProject removes a project record by ID.
func (r *SQLiteRepository) DeleteProject(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

// ── Environments ──────────────────────────────────────────────────────────────

// ListEnvironments returns all environments belonging to a project, value types (not pointers).
func (r *SQLiteRepository) ListEnvironments(_ context.Context, projectID string) ([]EnvironmentConfig, error) {
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

	var envs []EnvironmentConfig
	for rows.Next() {
		var env EnvironmentConfig
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}
		env.IsDefault = isDefault == 1
		envs = append(envs, env)
	}
	return envs, rows.Err()
}

// CreateEnvironment inserts a new environment record.
func (r *SQLiteRepository) CreateEnvironment(_ context.Context, env *EnvironmentConfig) error {
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

// DeleteEnvironment removes an environment by ID.
func (r *SQLiteRepository) DeleteEnvironment(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM environments WHERE id = ?`, id)
	return err
}

// getEnvironment is an internal helper that retrieves a single environment by ID (no context, no lock).
func (r *SQLiteRepository) getEnvironment(id string) (*EnvironmentConfig, error) {
	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, id,
	)
	var env EnvironmentConfig
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

// getDefaultEnvironment retrieves the default environment for a project (internal helper).
func (r *SQLiteRepository) getDefaultEnvironment(projectID string) (*EnvironmentConfig, error) {
	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? AND is_default = 1 LIMIT 1`,
		projectID,
	)
	var env EnvironmentConfig
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		// Fall back to the earliest created environment
		fallback := r.db.QueryRow(
			`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE project_id = ? ORDER BY created_at ASC LIMIT 1`,
			projectID,
		)
		err = fallback.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	env.IsDefault = isDefault == 1
	return &env, nil
}

// listAllEnvironments returns every environment across all projects (internal helper for canvas).
func (r *SQLiteRepository) listAllEnvironments() ([]*EnvironmentConfig, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments ORDER BY is_default DESC, created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []*EnvironmentConfig
	for rows.Next() {
		var env EnvironmentConfig
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, err
		}
		env.IsDefault = isDefault == 1
		envs = append(envs, &env)
	}
	return envs, rows.Err()
}

// ── Domains ───────────────────────────────────────────────────────────────────

// ListDomains returns all custom domains for the given project.
func (r *SQLiteRepository) ListDomains(_ context.Context, projectID string) ([]DomainConfig, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, domain_name, redirect_to, ssl_cert_status, path_prefix, created_at, updated_at FROM domains WHERE project_id = ? ORDER BY domain_name ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []DomainConfig
	for rows.Next() {
		var d DomainConfig
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.DomainName, &d.RedirectTo, &d.SSLCertStatus, &d.PathPrefix, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

// AddDomain inserts a new custom domain record.
func (r *SQLiteRepository) AddDomain(_ context.Context, d *DomainConfig) error {
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

// DeleteDomain removes a custom domain by ID.
func (r *SQLiteRepository) DeleteDomain(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM domains WHERE id = ?`, id)
	return err
}

// ── Env Vars ──────────────────────────────────────────────────────────────────

// GetEnvVars retrieves and decrypts all environment variables for a project.
func (r *SQLiteRepository) GetEnvVars(_ context.Context, projectID string) (map[string]string, error) {
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
			continue // skip undecryptable values rather than failing entirely
		}
		envs[key] = plaintext
	}
	return envs, rows.Err()
}

// SetEnvVar encrypts and upserts a single environment variable.
func (r *SQLiteRepository) SetEnvVar(_ context.Context, projectID, key, plaintextValue string) error {
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

// ── Canvas ────────────────────────────────────────────────────────────────────

// ListProjectCanvasSummaries returns dashboard summaries for every project without N+1 queries.
func (r *SQLiteRepository) ListProjectCanvasSummaries(ctx context.Context) ([]ProjectCanvasSummary, error) {
	projects, err := r.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	allEnvs, err := r.listAllEnvironments()
	if err != nil {
		return nil, fmt.Errorf("failed to list all environments: %w", err)
	}
	allApps, err := r.listAllAppServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services: %w", err)
	}
	allDbs, err := r.listAllDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list all databases: %w", err)
	}
	allStorage, err := r.listAllStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to list all storage: %w", err)
	}

	// Group by project
	envsByProject := make(map[string][]*EnvironmentConfig)
	for _, e := range allEnvs {
		envsByProject[e.ProjectID] = append(envsByProject[e.ProjectID], e)
	}
	appsByProject := make(map[string][]*types.AppServiceConfig)
	for _, a := range allApps {
		appsByProject[a.ProjectID] = append(appsByProject[a.ProjectID], a)
	}
	dbsByProject := make(map[string][]types.DatabaseConfig)
	for _, d := range allDbs {
		dbsByProject[d.ProjectID] = append(dbsByProject[d.ProjectID], d)
	}
	storageByProject := make(map[string][]types.StorageConfig)
	for _, s := range allStorage {
		storageByProject[s.ProjectID] = append(storageByProject[s.ProjectID], s)
	}

	var summaries []ProjectCanvasSummary
	for _, project := range projects {
		envs := envsByProject[project.ID]
		apps := appsByProject[project.ID]
		dbs := dbsByProject[project.ID]
		storage := storageByProject[project.ID]

		var defaultEnv *EnvironmentConfig
		if len(envs) > 0 {
			for _, e := range envs {
				if e.IsDefault {
					defaultEnv = e
					break
				}
			}
			if defaultEnv == nil {
				defaultEnv = envs[0]
			}
		}

		summary := ProjectCanvasSummary{
			ProjectConfig:      project,
			EnvironmentsCount:  len(envs),
			AppsCount:          len(apps),
			DatabasesCount:     len(dbs),
			StorageCount:       len(storage),
			TotalServices:      len(apps) + len(dbs) + len(storage),
			DefaultEnvironment: defaultEnv,
			ServiceIcons:       make([]string, 0),
		}

		onlineCount := 0
		for _, app := range apps {
			if app.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, "github")
		}
		for _, db := range dbs {
			if db.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, db.Engine)
		}
		for _, st := range storage {
			if st.Status == "running" {
				onlineCount++
			}
			summary.ServiceIcons = append(summary.ServiceIcons, st.Type)
		}
		summary.OnlineServices = onlineCount
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

// GetProjectCanvasSummary returns a canvas summary for a single project.
func (r *SQLiteRepository) GetProjectCanvasSummary(ctx context.Context, id string) (*ProjectCanvasSummary, error) {
	project, err := r.GetProject(ctx, id)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	envs, _ := r.ListEnvironments(ctx, id)
	apps, _ := r.listAppServicesByProject(id)
	dbs, _ := r.listDatabasesByProject(id)
	storage, _ := r.listStorageByProject(id)

	defaultEnv, _ := r.getDefaultEnvironment(id)

	summary := &ProjectCanvasSummary{
		ProjectConfig:      *project,
		EnvironmentsCount:  len(envs),
		AppsCount:          len(apps),
		DatabasesCount:     len(dbs),
		StorageCount:       len(storage),
		TotalServices:      len(apps) + len(dbs) + len(storage),
		DefaultEnvironment: defaultEnv,
		ServiceIcons:       make([]string, 0),
	}

	onlineCount := 0
	for _, app := range apps {
		if app.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, "github")
	}
	for _, db := range dbs {
		if db.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, db.Engine)
	}
	for _, st := range storage {
		if st.Status == "running" {
			onlineCount++
		}
		summary.ServiceIcons = append(summary.ServiceIcons, st.Type)
	}
	summary.OnlineServices = onlineCount
	return summary, nil
}

// GetEnvironmentCanvas retrieves all apps, databases, and storage for a given environment.
func (r *SQLiteRepository) GetEnvironmentCanvas(_ context.Context, environmentID string) (*EnvironmentCanvas, error) {
	env, err := r.getEnvironment(environmentID)
	if err != nil || env == nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	apps, _ := r.listAppServicesByEnvironment(environmentID)
	dbs, _ := r.listDatabasesByEnvironment(environmentID)
	storage, _ := r.listStorageByEnvironment(environmentID)

	var dbsPtrs []*types.DatabaseConfig
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	var storagePtrs []*types.StorageConfig
	for i := range storage {
		storagePtrs = append(storagePtrs, &storage[i])
	}

	return &EnvironmentCanvas{
		Environment: env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}, nil
}

// ── Internal read helpers (delegating to store tables) ────────────────────────

func (r *SQLiteRepository) listAllAppServices() ([]*types.AppServiceConfig, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(icon,''), COALESCE(repository_url,''), COALESCE(branch,''), COALESCE(root_directory,''), COALESCE(build_command,''), COALESCE(start_command,''), COALESCE(dockerfile_path,''), internal_port, COALESCE(domain,''), env_vars_count, auto_deploy_webhook, COALESCE(git_repo_full_name,''), wait_for_ci, auto_deploy_branch, COALESCE(public_networking_domain,''), COALESCE(private_networking_internal,''), enable_outbound_ipv6, cpu_request, memory_limit_mb, replicas, COALESCE(restart_policy,''), teardown_timeout, serverless, COALESCE(cron_schedule,''), COALESCE(health_check_path,''), status, COALESCE(container_id,''), created_at, updated_at FROM app_services ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listAppServicesByProject(projectID string) ([]*types.AppServiceConfig, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(icon,''), COALESCE(repository_url,''), COALESCE(branch,''), COALESCE(root_directory,''), COALESCE(build_command,''), COALESCE(start_command,''), COALESCE(dockerfile_path,''), internal_port, COALESCE(domain,''), env_vars_count, auto_deploy_webhook, COALESCE(git_repo_full_name,''), wait_for_ci, auto_deploy_branch, COALESCE(public_networking_domain,''), COALESCE(private_networking_internal,''), enable_outbound_ipv6, cpu_request, memory_limit_mb, replicas, COALESCE(restart_policy,''), teardown_timeout, serverless, COALESCE(cron_schedule,''), COALESCE(health_check_path,''), status, COALESCE(container_id,''), created_at, updated_at FROM app_services WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listAppServicesByEnvironment(environmentID string) ([]*types.AppServiceConfig, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(icon,''), COALESCE(repository_url,''), COALESCE(branch,''), COALESCE(root_directory,''), COALESCE(build_command,''), COALESCE(start_command,''), COALESCE(dockerfile_path,''), internal_port, COALESCE(domain,''), env_vars_count, auto_deploy_webhook, COALESCE(git_repo_full_name,''), wait_for_ci, auto_deploy_branch, COALESCE(public_networking_domain,''), COALESCE(private_networking_internal,''), enable_outbound_ipv6, cpu_request, memory_limit_mb, replicas, COALESCE(restart_policy,''), teardown_timeout, serverless, COALESCE(cron_schedule,''), COALESCE(health_check_path,''), status, COALESCE(container_id,''), created_at, updated_at FROM app_services WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanAppServices(query string, args ...any) ([]*types.AppServiceConfig, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*types.AppServiceConfig
	for rows.Next() {
		var a types.AppServiceConfig
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.EnvironmentID, &a.Name, &a.Icon, &a.RepositoryURL, &a.Branch,
			&a.RootDirectory, &a.BuildCommand, &a.StartCommand, &a.DockerfilePath, &a.InternalPort,
			&a.Domain, &a.EnvVarsCount, &a.AutoDeployWebhook, &a.GitRepoFullName, &a.WaitForCI,
			&a.AutoDeployBranch, &a.PublicNetworkingDomain, &a.PrivateNetworkingInternal,
			&a.EnableOutboundIPv6, &a.CPURequest, &a.MemoryLimitMB, &a.Replicas,
			&a.RestartPolicy, &a.TeardownTimeout, &a.Serverless, &a.CronSchedule,
			&a.HealthCheckPath, &a.Status, &a.ContainerID, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		apps = append(apps, &a)
	}
	return apps, rows.Err()
}

func (r *SQLiteRepository) listAllDatabases() ([]types.DatabaseConfig, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listDatabasesByProject(projectID string) ([]types.DatabaseConfig, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listDatabasesByEnvironment(environmentID string) ([]types.DatabaseConfig, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanDatabases(query string, args ...any) ([]types.DatabaseConfig, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []types.DatabaseConfig
	for rows.Next() {
		var d types.DatabaseConfig
		if err := rows.Scan(
			&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version, &d.Port,
			&d.Username, &d.DatabaseName, &d.VolumePath, &d.ContainerID, &d.Status,
			&d.InternalDNS, &d.ExternalDNS, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		dbs = append(dbs, d)
	}
	return dbs, rows.Err()
}

func (r *SQLiteRepository) listAllStorage() ([]types.StorageConfig, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listStorageByProject(projectID string) ([]types.StorageConfig, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listStorageByEnvironment(environmentID string) ([]types.StorageConfig, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanStorage(query string, args ...any) ([]types.StorageConfig, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.StorageConfig
	for rows.Next() {
		var s types.StorageConfig
		if err := rows.Scan(
			&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Name, &s.Type, &s.APIPort, &s.ConsolePort,
			&s.AccessKey, &s.BucketName, &s.VolumePath, &s.ContainerID, &s.Status,
			&s.InternalDNS, &s.ExternalDNS, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}
