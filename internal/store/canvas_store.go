package store

import (
	"fmt"

	"vessel.dev/vessel/internal/types"
)

// GetProjectCanvasSummary calculates aggregated counts and status icons for a single project canvas.
func (s *Store) GetProjectCanvasSummary(projectID string) (*types.ProjectCanvasSummary, error) {
	project, err := s.GetProject(projectID)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	envs, _ := s.ListEnvironments(projectID)
	apps, _ := s.ListAppServicesByProject(projectID)
	dbs, _ := s.ListDatabasesByProject(projectID)
	storage, _ := s.ListStorageByProject(projectID)

	defaultEnv, _ := s.GetDefaultEnvironment(projectID)

	summary := &types.ProjectCanvasSummary{
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

// ListProjectCanvasSummaries returns canvas summaries for all projects on the dashboard without N+1 queries.
func (s *Store) ListProjectCanvasSummaries() ([]*types.ProjectCanvasSummary, error) {
	projects, err := s.ListProjects()
	if err != nil {
		return nil, err
	}

	allEnvs, err := s.ListAllEnvironments()
	if err != nil {
		return nil, fmt.Errorf("failed to list all environments for canvas summaries: %w", err)
	}
	allApps, err := s.ListAllAppServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list all app services for canvas summaries: %w", err)
	}
	allDbs, err := s.ListAllDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list all databases for canvas summaries: %w", err)
	}
	allStorage, err := s.ListAllStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to list all storage for canvas summaries: %w", err)
	}

	envsByProject := make(map[string][]*types.EnvironmentConfig)
	for _, env := range allEnvs {
		envsByProject[env.ProjectID] = append(envsByProject[env.ProjectID], env)
	}

	appsByProject := make(map[string][]*types.AppServiceConfig)
	for _, app := range allApps {
		appsByProject[app.ProjectID] = append(appsByProject[app.ProjectID], app)
	}

	dbsByProject := make(map[string][]types.DatabaseConfig)
	for _, db := range allDbs {
		dbsByProject[db.ProjectID] = append(dbsByProject[db.ProjectID], db)
	}

	storageByProject := make(map[string][]types.StorageConfig)
	for _, st := range allStorage {
		storageByProject[st.ProjectID] = append(storageByProject[st.ProjectID], st)
	}

	var summaries []*types.ProjectCanvasSummary
	for _, project := range projects {
		envs := envsByProject[project.ID]
		apps := appsByProject[project.ID]
		dbs := dbsByProject[project.ID]
		storage := storageByProject[project.ID]

		var defaultEnv *types.EnvironmentConfig
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

		summary := &types.ProjectCanvasSummary{
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

// GetEnvironmentCanvas retrieves all Git applications, databases, and storage buckets inside a specific environment.
func (s *Store) GetEnvironmentCanvas(environmentID string) (*types.EnvironmentCanvas, error) {
	env, err := s.GetEnvironment(environmentID)
	if err != nil || env == nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	apps, _ := s.ListAppServicesByEnvironment(environmentID)
	dbs, _ := s.ListDatabasesByEnvironment(environmentID)
	storage, _ := s.ListStorageByEnvironment(environmentID)

	var dbsPtrs []*types.DatabaseConfig
	for i := range dbs {
		dbsPtrs = append(dbsPtrs, &dbs[i])
	}
	var storagePtrs []*types.StorageConfig
	for i := range storage {
		storagePtrs = append(storagePtrs, &storage[i])
	}

	canvas := &types.EnvironmentCanvas{
		Environment: env,
		Apps:        apps,
		Databases:   dbsPtrs,
		Storage:     storagePtrs,
	}
	return canvas, nil
}
