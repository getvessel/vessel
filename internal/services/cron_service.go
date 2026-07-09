package services

import (
	"context"
	"errors"
	"fmt"

	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/project"
)

type JobStore interface {
	CreateJob(job *job.Job) error
	GetJob(id string) (*job.Job, error)
	ListJobs() ([]job.Job, error)
	ListJobsByProject(projectID string) ([]job.Job, error)
	DeleteJob(id string) error
}

type ProjectStore interface {
	GetProject(id string) (*project.ProjectConfig, error)
}

type CronService struct {
	jobs        JobStore
	projects    ProjectStore
	cronManager *orchestrator.CronManager
}

func NewCronService(js JobStore, ps ProjectStore, cm *orchestrator.CronManager) *CronService {
	return &CronService{
		jobs:        js,
		projects:    ps,
		cronManager: cm,
	}
}

// CreateJob persists a new job and immediately registers its cron timer if active.
func (cs *CronService) CreateJob(j *job.Job) error {
	if j.ProjectID == "" {
		return errors.New("projectId is required when creating a scheduled job")
	}
	if j.Schedule == "" {
		return errors.New("schedule cron expression is required")
	}
	if j.Command == "" {
		return errors.New("command is required")
	}

	project, err := cs.projects.GetProject(j.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to verify project existence: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project with ID %s not found", j.ProjectID)
	}

	if err := cs.jobs.CreateJob(j); err != nil {
		return err
	}

	return cs.cronManager.RegisterJob(j)
}

// GetJob retrieves a job configuration by ID.
func (cs *CronService) GetJob(id string) (*job.Job, error) {
	return cs.jobs.GetJob(id)
}

// ListJobs retrieves all scheduled jobs, optionally filtered by projectID.
func (cs *CronService) ListJobs(projectID string) ([]job.Job, error) {
	if projectID != "" {
		return cs.jobs.ListJobsByProject(projectID)
	}
	return cs.jobs.ListJobs()
}

// TriggerJobImmediately executes a scheduled job synchronously on demand.
func (cs *CronService) TriggerJobImmediately(ctx context.Context, jobID string) (string, error) {
	return cs.cronManager.ExecuteJob(ctx, jobID)
}

// DeleteJob permanently removes a job from persistence and halts its active cron timer.
func (cs *CronService) DeleteJob(id string) error {
	cs.cronManager.UnregisterJob(id)
	return cs.jobs.DeleteJob(id)
}
