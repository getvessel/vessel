package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
)

// CreateJob inserts a new scheduled cron job record into the store.
func (s *Store) CreateJob(job *types.JobConfig) error {
	if job.ID == "" {
		job.ID = uuid.NewString()
	}
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	if job.Status == "" {
		job.Status = "active"
	}

	_, err := s.db.Exec(`INSERT INTO jobs (
		id, project_id, name, schedule, command, status, last_run_at, last_output, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		job.ID, job.ProjectID, job.Name, job.Schedule, job.Command, job.Status, job.LastRunAt, job.LastOutput, job.CreatedAt, job.UpdatedAt)
	return err
}

// GetJob retrieves a specific scheduled job by ID.
func (s *Store) GetJob(id string) (*types.JobConfig, error) {
	var job types.JobConfig
	var lastRunAt sql.NullTime

	err := s.db.QueryRow(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at
		FROM jobs WHERE id = ?`, id).Scan(
		&job.ID, &job.ProjectID, &job.Name, &job.Schedule, &job.Command, &job.Status, &lastRunAt, &job.LastOutput, &job.CreatedAt, &job.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lastRunAt.Valid {
		job.LastRunAt = &lastRunAt.Time
	}
	return &job, nil
}

// ListJobs retrieves all scheduled background jobs across all projects.
func (s *Store) ListJobs() ([]types.JobConfig, error) {
	rows, err := s.db.Query(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at FROM jobs ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []types.JobConfig
	for rows.Next() {
		var job types.JobConfig
		var lastRunAt sql.NullTime
		if err := rows.Scan(&job.ID, &job.ProjectID, &job.Name, &job.Schedule, &job.Command, &job.Status, &lastRunAt, &job.LastOutput, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		if lastRunAt.Valid {
			job.LastRunAt = &lastRunAt.Time
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// ListJobsByProject retrieves all scheduled jobs associated with a specific project container.
func (s *Store) ListJobsByProject(projectID string) ([]types.JobConfig, error) {
	rows, err := s.db.Query(`SELECT id, project_id, name, schedule, command, status, last_run_at, COALESCE(last_output, ''), created_at, updated_at FROM jobs WHERE project_id = ? ORDER BY created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []types.JobConfig
	for rows.Next() {
		var job types.JobConfig
		var lastRunAt sql.NullTime
		if err := rows.Scan(&job.ID, &job.ProjectID, &job.Name, &job.Schedule, &job.Command, &job.Status, &lastRunAt, &job.LastOutput, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		if lastRunAt.Valid {
			job.LastRunAt = &lastRunAt.Time
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// UpdateJobStatusAndOutput updates the execution state, completion timestamp, and captured console logs of a job.
func (s *Store) UpdateJobStatusAndOutput(id string, status string, lastRunAt *time.Time, output string) error {
	now := time.Now()
	_, err := s.db.Exec(`UPDATE jobs SET status = ?, last_run_at = ?, last_output = ?, updated_at = ? WHERE id = ?`,
		status, lastRunAt, output, now, id)
	return err
}

// DeleteJob permanently removes a scheduled job from the store.
func (s *Store) DeleteJob(id string) error {
	_, err := s.db.Exec(`DELETE FROM jobs WHERE id = ?`, id)
	return err
}
