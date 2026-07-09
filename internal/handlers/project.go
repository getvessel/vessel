package handlers

import (
	"github.com/labstack/echo/v4"

	"context"
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type ProxyReloader interface {
	Reload(ctx context.Context) error
}

type ProjectHandler struct {
	projectService *services.ProjectService
	proxy          ProxyReloader
}

func NewProjectHandler(s *services.ProjectService, p ProxyReloader) *ProjectHandler {
	return &ProjectHandler{projectService: s, proxy: p}
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	projects, err := h.projectService.ListProjects(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	p, err := h.projectService.CreateProjectFromRequest(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	WriteJSON(w, http.StatusCreated, p)
}

func (h *ProjectHandler) GetProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	p, err := h.projectService.GetProject(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "project not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	if err := h.projectService.DeleteProject(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
