package store

import (
	"time"

	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/storage"
	"vessel.dev/vessel/internal/types"
)

// ── Proxy adapters ────────────────────────────────────────────────────────────

type ProxyStore struct {
	*Store
}

func NewProxyStore(s *Store) *ProxyStore {
	return &ProxyStore{Store: s}
}

func (s *ProxyStore) ListProjects() ([]project.ProjectConfig, error) {
	items, err := s.Store.ListProjects()
	if err != nil {
		return nil, err
	}
	result := make([]project.ProjectConfig, len(items))
	for i, p := range items {
		result[i] = project.ProjectConfig{
			ID: p.ID, WorkspaceID: p.WorkspaceID, TeamID: p.TeamID,
			Name: p.Name, Description: p.Description,
			CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
		}
	}
	return result, nil
}

func (s *ProxyStore) ListAllAppServices() ([]*service.AppService, error) {
	items, err := s.Store.ListAllAppServices()
	if err != nil {
		return nil, err
	}
	result := make([]*service.AppService, len(items))
	for i, a := range items {
		if a == nil {
			continue
		}
		result[i] = &service.AppService{
			ID: a.ID, ProjectID: a.ProjectID, EnvironmentID: a.EnvironmentID,
			Name: a.Name, RepositoryURL: a.RepositoryURL, Branch: a.Branch,
			InternalPort: a.InternalPort, Domain: a.Domain,
			ContainerID: a.ContainerID, Status: a.Status,
			CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
		}
	}
	return result, nil
}

func (s *ProxyStore) ListAllDomains() ([]domain.Config, error) {
	items, err := s.Store.ListAllDomains()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Config, len(items))
	for i, d := range items {
		result[i] = domain.Config{
			ID: d.ID, ProjectID: d.ProjectID, DomainName: d.DomainName,
			RedirectTo: d.RedirectTo, SSLCertStatus: d.SSLCertStatus,
			PathPrefix: d.PathPrefix, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		}
	}
	return result, nil
}

// ── Cron/Job adapters ─────────────────────────────────────────────────────────

type CronJobStore struct {
	*Store
}

func NewCronJobStore(s *Store) *CronJobStore {
	return &CronJobStore{Store: s}
}

func (s *CronJobStore) CreateJob(j *job.Job) error {
	return s.Store.CreateJob(toTypesJob(j))
}

func (s *CronJobStore) GetJob(id string) (*job.Job, error) {
	cfg, err := s.Store.GetJob(id)
	if err != nil || cfg == nil {
		return nil, err
	}
	return fromTypesJob(cfg), nil
}

func (s *CronJobStore) ListJobs() ([]job.Job, error) {
	items, err := s.Store.ListJobs()
	if err != nil {
		return nil, err
	}
	result := make([]job.Job, len(items))
	for i, cfg := range items {
		result[i] = *fromTypesJob(&cfg)
	}
	return result, nil
}

func (s *CronJobStore) ListJobsByProject(projectID string) ([]job.Job, error) {
	items, err := s.Store.ListJobsByProject(projectID)
	if err != nil {
		return nil, err
	}
	result := make([]job.Job, len(items))
	for i, cfg := range items {
		result[i] = *fromTypesJob(&cfg)
	}
	return result, nil
}

func (s *CronJobStore) DeleteJob(id string) error {
	return s.Store.DeleteJob(id)
}

type CronProjectStore struct {
	*Store
}

func NewCronProjectStore(s *Store) *CronProjectStore {
	return &CronProjectStore{Store: s}
}

func (s *CronProjectStore) GetProject(id string) (*project.ProjectConfig, error) {
	cfg, err := s.Store.GetProject(id)
	if err != nil || cfg == nil {
		return nil, err
	}
	return &project.ProjectConfig{
		ID: cfg.ID, WorkspaceID: cfg.WorkspaceID, TeamID: cfg.TeamID,
		Name: cfg.Name, Description: cfg.Description,
		CreatedAt: cfg.CreatedAt, UpdatedAt: cfg.UpdatedAt,
	}, nil
}

// ── Service Linker adapters ───────────────────────────────────────────────────

type LinkerDatabaseStore struct {
	*Store
}

func NewLinkerDatabaseStore(s *Store) *LinkerDatabaseStore {
	return &LinkerDatabaseStore{Store: s}
}

func (s *LinkerDatabaseStore) ListDatabasesByProject(projectID string) ([]database.Database, error) {
	items, err := s.Store.ListDatabasesByProject(projectID)
	if err != nil {
		return nil, err
	}
	result := make([]database.Database, len(items))
	for i, d := range items {
		result[i] = database.Database{
			ID: d.ID, ProjectID: d.ProjectID, EnvironmentID: d.EnvironmentID,
			Name: d.Name, Engine: d.Engine, Version: d.Version,
			Port: d.Port, Username: d.Username, Password: d.Password,
			DatabaseName: d.DatabaseName, VolumePath: d.VolumePath,
			ContainerID: d.ContainerID, Status: d.Status,
			InternalDNS: d.InternalDNS, ExternalDNS: d.ExternalDNS,
			CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		}
	}
	return result, nil
}

type LinkerStorageStore struct {
	*Store
}

func NewLinkerStorageStore(s *Store) *LinkerStorageStore {
	return &LinkerStorageStore{Store: s}
}

func (s *LinkerStorageStore) ListStorageByProject(projectID string) ([]storage.Storage, error) {
	items, err := s.Store.ListStorageByProject(projectID)
	if err != nil {
		return nil, err
	}
	result := make([]storage.Storage, len(items))
	for i, st := range items {
		result[i] = storage.Storage{
			ID: st.ID, ProjectID: st.ProjectID, EnvironmentID: st.EnvironmentID,
			Name: st.Name, Type: st.Type, APIPort: st.APIPort, ConsolePort: st.ConsolePort,
			AccessKey: st.AccessKey, SecretKey: st.SecretKey, BucketName: st.BucketName,
			VolumePath: st.VolumePath, ContainerID: st.ContainerID, Status: st.Status,
			InternalDNS: st.InternalDNS, ExternalDNS: st.ExternalDNS,
			CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt,
		}
	}
	return result, nil
}

// ── Type conversion helpers ───────────────────────────────────────────────────

func toTypesJob(j *job.Job) *types.JobConfig {
	return &types.JobConfig{
		ID: j.ID, ProjectID: j.ProjectID, Name: j.Name,
		Schedule: j.Schedule, Command: j.Command, Status: j.Status,
		LastRunAt: j.LastRunAt, LastOutput: j.LastOutput,
		CreatedAt: j.CreatedAt, UpdatedAt: j.UpdatedAt,
	}
}

func fromTypesJob(cfg *types.JobConfig) *job.Job {
	var lastRunAt *time.Time
	if cfg.LastRunAt != nil {
		lastRunAt = cfg.LastRunAt
	}
	return &job.Job{
		ID: cfg.ID, ProjectID: cfg.ProjectID, Name: cfg.Name,
		Schedule: cfg.Schedule, Command: cfg.Command, Status: cfg.Status,
		LastRunAt: lastRunAt, LastOutput: cfg.LastOutput,
		CreatedAt: cfg.CreatedAt, UpdatedAt: cfg.UpdatedAt,
	}
}
