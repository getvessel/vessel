package job

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) ListProjectJobs(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("projectId")
	jobs, err := h.repo.ListByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list jobs: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	job := &Job{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Schedule:  req.Schedule,
		Command:   req.Command,
	}

	if err := h.repo.Create(r.Context(), job); err != nil {
		writeError(w, http.StatusBadRequest, "failed to create job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, job)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve job: "+err.Error())
		return
	}
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	var req UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Schedule != nil {
		existing.Schedule = *req.Schedule
	}
	if req.Command != nil {
		existing.Command = *req.Command
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := h.repo.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "job deleted successfully"})
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, "running", nil, ""); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to trigger job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "triggered", "jobId": id})
}

func (h *Handler) Pause(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, "paused", nil, ""); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to pause job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "paused", "jobId": id})
}

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, "active", nil, ""); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to resume job: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resumed", "jobId": id})
}
