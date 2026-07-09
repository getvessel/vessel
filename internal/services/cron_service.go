package services

import (
	"context"
	"errors"
	"fmt"

	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/engine"
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
	cronManager *engine.CronManager
}

func NewCronService(js JobStore, ps ProjectStore, cm *engine.CronManager) *CronService {
	return &CronService{
		jobs:        js,
		projects:    ps,
		cronManager: cm,
	}
}

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

func (cs *CronService) GetJob(id string) (*job.Job, error) {
	return cs.jobs.GetJob(id)
}

func (cs *CronService) ListJobs(projectID string) ([]job.Job, error) {
	if projectID != "" {
		return cs.jobs.ListJobsByProject(projectID)
	}
	return cs.jobs.ListJobs()
}

func (cs *CronService) TriggerJobImmediately(ctx context.Context, jobID string) (string, error) {
	return cs.cronManager.ExecuteJob(ctx, jobID)
}

func (cs *CronService) DeleteJob(id string) error {
	cs.cronManager.UnregisterJob(id)
	return cs.jobs.DeleteJob(id)
}
