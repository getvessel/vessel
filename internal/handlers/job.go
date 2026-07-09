package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *JobHandler) ListProjectJobs(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId query parameter")
		return nil
	}
	jobs, err := h.jobService.ListJobsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, jobs)
}

func (h *JobHandler) Create(c echo.Context) error {
	var j models.Job
	if err := c.Bind(&j); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	created, err := h.jobService.CreateJob(r.Context(), &j)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *JobHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	j, err := h.jobService.GetJob(r.Context(), id)
	if err != nil || j == nil {
		WriteError(w, http.StatusNotFound, "job not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, j)
}

func (h *JobHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	if err := h.jobService.DeleteJob(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) Run(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	out, err := h.jobService.ExecuteJob(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "executed", "output": out})
}
