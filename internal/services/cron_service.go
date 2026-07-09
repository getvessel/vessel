package services

import (
	"context"
	"errors"
	"fmt"

	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

type CronService struct {
	store       *store.Store
	cronManager *orchestrator.CronManager
}

func NewCronService(s *store.Store, cm *orchestrator.CronManager) *CronService {
	return &CronService{
		store:       s,
		cronManager: cm,
	}
}

// CreateJob persists a new job and immediately registers its cron timer if active.
func (cs *CronService) CreateJob(job *types.JobConfig) error {
	if job.ProjectID == "" {
		return errors.New("projectId is required when creating a scheduled job")
	}
	if job.Schedule == "" {
		return errors.New("schedule cron expression is required")
	}
	if job.Command == "" {
		return errors.New("command is required")
	}

	project, err := cs.store.GetProject(job.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to verify project existence: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project with ID %s not found", job.ProjectID)
	}

	if err := cs.store.CreateJob(job); err != nil {
		return err
	}

	return cs.cronManager.RegisterJob(job)
}

// GetJob retrieves a job configuration by ID.
func (cs *CronService) GetJob(id string) (*types.JobConfig, error) {
	return cs.store.GetJob(id)
}

// ListJobs retrieves all scheduled jobs, optionally filtered by projectID.
func (cs *CronService) ListJobs(projectID string) ([]types.JobConfig, error) {
	if projectID != "" {
		return cs.store.ListJobsByProject(projectID)
	}
	return cs.store.ListJobs()
}

// TriggerJobImmediately executes a scheduled job synchronously on demand.
func (cs *CronService) TriggerJobImmediately(ctx context.Context, jobID string) (string, error) {
	return cs.cronManager.ExecuteJob(ctx, jobID)
}

// DeleteJob permanently removes a job from persistence and halts its active cron timer.
func (cs *CronService) DeleteJob(id string) error {
	cs.cronManager.UnregisterJob(id)
	return cs.store.DeleteJob(id)
}
