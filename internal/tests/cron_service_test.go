package tests

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestCronServiceCreateAndListJobs(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_cron_service_test_create")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	cm := orchestrator.NewCronManager(nil, s)
	cs := services.NewCronService(s, cm)

	err = cs.CreateJob(&types.JobConfig{
		Name:     "Test Job",
		Schedule: "0 * * * *",
		Command:  "echo hello",
	})
	if err == nil || !strings.Contains(err.Error(), "projectId is required") {
		t.Fatalf("expected projectId required error, got: %v", err)
	}

	err = cs.CreateJob(&types.JobConfig{
		ProjectID: "proj-1",
		Name:      "Test Job",
		Command:   "echo hello",
	})
	if err == nil || !strings.Contains(err.Error(), "schedule cron expression is required") {
		t.Fatalf("expected schedule required error, got: %v", err)
	}

	err = cs.CreateJob(&types.JobConfig{
		ProjectID: "proj-1",
		Name:      "Test Job",
		Schedule:  "0 * * * *",
	})
	if err == nil || !strings.Contains(err.Error(), "command is required") {
		t.Fatalf("expected command required error, got: %v", err)
	}

	err = cs.CreateJob(&types.JobConfig{
		ProjectID: "non-existent-proj",
		Name:      "Test Job",
		Schedule:  "0 * * * *",
		Command:   "echo hello",
	})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected project not found error, got: %v", err)
	}

	proj := &types.ProjectConfig{
		Name: "Cron Project",
	}
	if err := s.CreateProject(proj); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	job := &types.JobConfig{
		ProjectID: proj.ID,
		Name:      "Valid Job",
		Schedule:  "*/5 * * * *",
		Command:   "echo valid",
		Status:    "active",
	}
	if err := cs.CreateJob(job); err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}
	if job.ID == "" {
		t.Fatal("expected job ID to be populated")
	}

	fetched, err := cs.GetJob(job.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if fetched == nil || fetched.Name != "Valid Job" {
		t.Fatalf("unexpected fetched job: %+v", fetched)
	}

	proj2 := &types.ProjectConfig{Name: "Second Project"}
	_ = s.CreateProject(proj2)
	job2 := &types.JobConfig{
		ProjectID: proj2.ID,
		Name:      "Job 2",
		Schedule:  "0 0 * * *",
		Command:   "echo job2",
		Status:    "inactive",
	}
	_ = cs.CreateJob(job2)

	jobsProj1, err := cs.ListJobs(proj.ID)
	if err != nil {
		t.Fatalf("ListJobs(proj.ID) failed: %v", err)
	}
	if len(jobsProj1) != 1 || jobsProj1[0].ID != job.ID {
		t.Fatalf("expected 1 job for proj 1, got %d", len(jobsProj1))
	}

	allJobs, err := cs.ListJobs("")
	if err != nil {
		t.Fatalf("ListJobs(\"\") failed: %v", err)
	}
	if len(allJobs) != 2 {
		t.Fatalf("expected 2 jobs total, got %d", len(allJobs))
	}

	if err := cs.DeleteJob(job.ID); err != nil {
		t.Fatalf("DeleteJob failed: %v", err)
	}
	deleted, _ := cs.GetJob(job.ID)
	if deleted != nil {
		t.Fatal("expected job to be deleted")
	}
}

func TestCronServiceTriggerJobImmediately(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_cron_service_test_trigger")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "vessel.db")
	s, err := store.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	mockClient, _ := client.NewClientWithOpts(client.WithHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"message": "No such container"}`)),
			}, nil
		}),
	}))

	cm := orchestrator.NewCronManager(mockClient, s)
	cs := services.NewCronService(s, cm)

	_, err = cs.TriggerJobImmediately(context.Background(), "non-existent")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected job not found error, got: %v", err)
	}

	proj := &types.ProjectConfig{Name: "Trigger Project"}
	_ = s.CreateProject(proj)
	job := &types.JobConfig{
		ProjectID: proj.ID,
		Name:      "Trigger Job",
		Schedule:  "0 * * * *",
		Command:   "echo run",
		Status:    "active",
	}
	_ = cs.CreateJob(job)

	_, err = cs.TriggerJobImmediately(context.Background(), job.ID)
	if err == nil || !strings.Contains(err.Error(), "stopped or not found") {
		t.Fatalf("expected stopped or not found error from trigger, got: %v", err)
	}

	updated, _ := cs.GetJob(job.ID)
	if updated.Status != "error" {
		t.Fatalf("expected status 'error', got %s", updated.Status)
	}
}
