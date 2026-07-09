package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/types"
)

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		projectID := strings.TrimSpace(r.URL.Query().Get("projectId"))
		jobs, err := s.cronService.ListJobs(projectID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list scheduled jobs: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, jobs)

	case http.MethodPost:
		var req types.JobConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body payload")
			return
		}

		if err := s.cronService.CreateJob(&req); err != nil {
			writeError(w, http.StatusBadRequest, "failed to create job: "+err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, req)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleJobDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	parts := strings.Split(path, "/")
	jobID := parts[0]
	if jobID == "" {
		writeError(w, http.StatusBadRequest, "job id parameter is required")
		return
	}

	if len(parts) > 1 && parts[1] == "trigger" && r.Method == http.MethodPost {
		output, err := s.cronService.TriggerJobImmediately(r.Context(), jobID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "job execution failed: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "success",
			"jobId":  jobID,
			"output": output,
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		job, err := s.cronService.GetJob(jobID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to retrieve job details: "+err.Error())
			return
		}
		if job == nil {
			writeError(w, http.StatusNotFound, "scheduled job not found")
			return
		}
		writeJSON(w, http.StatusOK, job)

	case http.MethodDelete:
		if err := s.cronService.DeleteJob(jobID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete scheduled job: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "job deleted successfully"})

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
