package canvas

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/storage"
)

// SQLiteRepository implements Repository against a SQLite database.
type SQLiteRepository struct {
	db           *sql.DB
	mu           sync.Mutex
	environments environment.Repository
}

// NewSQLiteRepository constructs a SQLiteRepository backed by the given db and environment repository.
func NewSQLiteRepository(db *sql.DB, envRepo environment.Repository) *SQLiteRepository {
	return &SQLiteRepository{db: db, environments: envRepo}
}

// ── Canvas read model ────────────────────────────────────────────────────────

// ListCanvasSummaries returns dashboard summaries for every project without N+1 queries.
func (r *SQLiteRepository) ListCanvasSummaries(ctx context.Context) ([]CanvasSummary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	projects, err := r.listAllProjects()
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
	appsByProject := make(map[string][]*service.AppService)
	for _, a := range allApps {
		appsByProject[a.ProjectID] = append(appsByProject[a.ProjectID], a)
	}
	dbsByProject := make(map[string][]database.Database)
	for _, d := range allDbs {
		dbsByProject[d.ProjectID] = append(dbsByProject[d.ProjectID], d)
	}
	storageByProject := make(map[string][]storage.Storage)
	for _, st := range allStorage {
		storageByProject[st.ProjectID] = append(storageByProject[st.ProjectID], st)
	}

	var summaries []CanvasSummary
	for _, project := range projects {
		envs := envsByProject[project.ID]
		apps := appsByProject[project.ID]
		dbs := dbsByProject[project.ID]
		storageItems := storageByProject[project.ID]

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
			ID:                 project.ID,
			WorkspaceID:        project.WorkspaceID,
			TeamID:             project.TeamID,
			Name:               project.Name,
			Description:        project.Description,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
			EnvironmentsCount:  len(envs),
			AppsCount:          len(apps),
			DatabasesCount:     len(dbs),
			StorageCount:       len(storageItems),
			TotalServices:      len(apps) + len(dbs) + len(storageItems),
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
		for _, st := range storageItems {
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
	r.mu.Lock()
	defer r.mu.Unlock()

	project, err := r.getProject(id)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	envs, _ := r.environments.ListByProject(ctx, id)
	apps, _ := r.listAppServicesByProject(id)
	dbs, _ := r.listDatabasesByProject(id)
	storageItems, _ := r.listStorageByProject(id)

	var defaultEnv *environment.Config
	if len(envs) > 0 {
		for _, e := range envs {
			e := e
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
		ID:                 project.ID,
		WorkspaceID:        project.WorkspaceID,
		TeamID:             project.TeamID,
		Name:               project.Name,
		Description:        project.Description,
		CreatedAt:          project.CreatedAt,
		UpdatedAt:          project.UpdatedAt,
		EnvironmentsCount:  len(envs),
		AppsCount:          len(apps),
		DatabasesCount:     len(dbs),
		StorageCount:       len(storageItems),
		TotalServices:      len(apps) + len(dbs) + len(storageItems),
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
	for _, st := range storageItems {
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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	storageItems, _ := r.listStorageByEnvironment(environmentID)

	var dbsPtrs []*database.Database
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	var storagePtrs []*storage.Storage
	for i := range storageItems {
		storagePtrs = append(storagePtrs, &storageItems[i])
	}

	return &EnvironmentCanvas{
		Environment: &env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}, nil
}

// ── Internal helpers ─────────────────────────────────────────────────────────

type projectRow struct {
	ID          string
	WorkspaceID string
	TeamID      string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *SQLiteRepository) listAllProjects() ([]projectRow, error) {
	rows, err := r.db.Query(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []projectRow
	for rows.Next() {
		var p projectRow
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *SQLiteRepository) getProject(id string) (*projectRow, error) {
	row := r.db.QueryRow(`SELECT id, COALESCE(workspace_id, ''), COALESCE(team_id,''), name, COALESCE(description,''), created_at, updated_at FROM projects WHERE id = ?`, id)
	var p projectRow
	err := row.Scan(&p.ID, &p.WorkspaceID, &p.TeamID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SQLiteRepository) listAllEnvironments(ctx context.Context) ([]*environment.Config, error) {
	envs, err := r.environments.ListByProject(ctx, "")
	if err == nil && len(envs) > 0 {
		var result []*environment.Config
		for i := range envs {
			result = append(result, &envs[i])
		}
		return result, nil
	}
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

func (r *SQLiteRepository) listAllAppServices() ([]*service.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listAppServicesByProject(projectID string) ([]*service.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listAppServicesByEnvironment(environmentID string) ([]*service.AppService, error) {
	return r.scanAppServices(`SELECT id, project_id, environment_id, name, COALESCE(repository_url,''), COALESCE(branch,''), internal_port, COALESCE(domain,''), COALESCE(container_id,''), status, created_at, updated_at FROM app_services WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanAppServices(query string, args ...any) ([]*service.AppService, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*service.AppService
	for rows.Next() {
		var a service.AppService
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.EnvironmentID, &a.Name,
			&a.RepositoryURL, &a.Branch, &a.InternalPort,
			&a.Domain, &a.ContainerID, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		apps = append(apps, &a)
	}
	return apps, rows.Err()
}

func (r *SQLiteRepository) listAllDatabases() ([]database.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listDatabasesByProject(projectID string) ([]database.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listDatabasesByEnvironment(environmentID string) ([]database.Database, error) {
	return r.scanDatabases(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, engine, version, port, username, database_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM databases WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanDatabases(query string, args ...any) ([]database.Database, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []database.Database
	for rows.Next() {
		var d database.Database
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

func (r *SQLiteRepository) listAllStorage() ([]storage.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage ORDER BY created_at DESC`)
}

func (r *SQLiteRepository) listStorageByProject(projectID string) ([]storage.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE project_id = ? ORDER BY created_at DESC`, projectID)
}

func (r *SQLiteRepository) listStorageByEnvironment(environmentID string) ([]storage.Storage, error) {
	return r.scanStorage(`SELECT id, COALESCE(project_id,''), COALESCE(environment_id,''), name, type, api_port, console_port, access_key, bucket_name, volume_path, COALESCE(container_id,''), status, COALESCE(internal_dns,''), COALESCE(external_dns,''), created_at, updated_at FROM storage WHERE environment_id = ? ORDER BY created_at DESC`, environmentID)
}

func (r *SQLiteRepository) scanStorage(query string, args ...any) ([]storage.Storage, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []storage.Storage
	for rows.Next() {
		var s storage.Storage
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
