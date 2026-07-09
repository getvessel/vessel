package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(s *services.JobService) *JobHandler {
	return &JobHandler{jobService: s}
}

func (h *JobHandler) ListProjectJobs(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return
	}
	jobs, err := h.jobService.ListJobsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, jobs)
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var j models.Job
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	created, err := h.jobService.CreateJob(r.Context(), &j)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	j, err := h.jobService.GetJob(r.Context(), id)
	if err != nil || j == nil {
		WriteError(w, http.StatusNotFound, "job not found")
		return
	}
	WriteJSON(w, http.StatusOK, j)
}

func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	if err := h.jobService.DeleteJob(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) Run(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	out, err := h.jobService.ExecuteJob(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "executed", "output": out})
}
