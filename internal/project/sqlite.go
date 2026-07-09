package project

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/types"
)

// SQLiteRepository implements Repository against a SQLite database.
type SQLiteRepository struct {
	db           *sql.DB
	environments environment.Repository
}

// NewSQLiteRepository constructs a SQLiteRepository backed by the given db and environment repository.
func NewSQLiteRepository(db *sql.DB, envRepo environment.Repository) *SQLiteRepository {
	return &SQLiteRepository{db: db, environments: envRepo}
}

// List retrieves all ProjectConfig records ordered by creation date descending.
func (r *SQLiteRepository) List(_ context.Context) ([]ProjectConfig, error) {
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

// Get retrieves a single ProjectConfig by its ID.
func (r *SQLiteRepository) Get(_ context.Context, id string) (*ProjectConfig, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p ProjectConfig
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create inserts a new project and creates its default production environment.
func (r *SQLiteRepository) Create(ctx context.Context, p *ProjectConfig) error {
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

	defaultEnv := &environment.Config{
		ProjectID: p.ID,
		Name:      "production",
		IsDefault: true,
	}
	return r.environments.Create(ctx, defaultEnv)
}

// Delete removes a project record by ID.
func (r *SQLiteRepository) Delete(_ context.Context, id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

// ── Canvas read model ────────────────────────────────────────────────────────

// ListCanvasSummaries returns dashboard summaries for every project without N+1 queries.
func (r *SQLiteRepository) ListCanvasSummaries(ctx context.Context) ([]CanvasSummary, error) {
	projects, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	allEnvs, err := r.listAllEnvironments(ctx)
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

	envsByProject := make(map[string][]*environment.Config)
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
	for _, st := range allStorage {
		storageByProject[st.ProjectID] = append(storageByProject[st.ProjectID], st)
	}

	var summaries []CanvasSummary
	for _, project := range projects {
		envs := envsByProject[project.ID]
		apps := appsByProject[project.ID]
		dbs := dbsByProject[project.ID]
		storage := storageByProject[project.ID]

		var defaultEnv *environment.Config
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

		summary := CanvasSummary{
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

// GetCanvasSummary returns a canvas summary for a single project.
func (r *SQLiteRepository) GetCanvasSummary(ctx context.Context, id string) (*CanvasSummary, error) {
	project, err := r.Get(ctx, id)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	envs, _ := r.environments.ListByProject(ctx, id)
	apps, _ := r.listAppServicesByProject(id)
	dbs, _ := r.listDatabasesByProject(id)
	storage, _ := r.listStorageByProject(id)

	var defaultEnv *environment.Config
	if len(envs) > 0 {
		for _, e := range envs {
			if e.IsDefault {
				defaultEnv = &e
				break
			}
		}
		if defaultEnv == nil {
			defaultEnv = &envs[0]
		}
	}

	summary := &CanvasSummary{
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
	row := r.db.QueryRow(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments WHERE id = ?`, environmentID,
	)
	var env environment.Config
	var isDefault int
	err := row.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("environment not found: %s", environmentID)
	}
	if err != nil {
		return nil, err
	}
	env.IsDefault = isDefault == 1

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
		Environment: &env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}, nil
}

// ── Internal read helpers (delegating to store tables) ────────────────────────

func (r *SQLiteRepository) listAllEnvironments(ctx context.Context) ([]*environment.Config, error) {
	envs, err := r.environments.ListByProject(ctx, "")
	if err == nil && len(envs) > 0 {
		var result []*environment.Config
		for i := range envs {
			result = append(result, &envs[i])
		}
		return result, nil
	}
	// Fallback: query directly so we don't require ListByProject to support empty projectID.
	rows, err := r.db.Query(
		`SELECT id, project_id, name, is_default, created_at, updated_at FROM environments ORDER BY is_default DESC, created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*environment.Config
	for rows.Next() {
		var env environment.Config
		var isDefault int
		if err := rows.Scan(&env.ID, &env.ProjectID, &env.Name, &isDefault, &env.CreatedAt, &env.UpdatedAt); err != nil {
			return nil, err
		}
		env.IsDefault = isDefault == 1
		result = append(result, &env)
	}
	return result, rows.Err()
}

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
