package handlers

import (
	"github.com/labstack/echo/v4"

	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type EnvironmentHandler struct {
	envService *services.EnvironmentService
}

func NewEnvironmentHandler(s *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{envService: s}
}

func (h *EnvironmentHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	envs, err := h.envService.ListByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, envs)
}

func (h *EnvironmentHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	var env models.EnvironmentConfig
	if err := c.Bind(&env); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	env.ProjectID = projectID
	if env.Name == "" {
		WriteError(w, http.StatusBadRequest, "environment name is required")
		return nil
	}
	created, err := h.envService.CreateEnvironment(r.Context(), &env)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *EnvironmentHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	if err := h.envService.DeleteEnvironment(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

type DomainHandler struct {
	envService *services.EnvironmentService
}

func NewDomainHandler(s *services.EnvironmentService) *DomainHandler {
	return &DomainHandler{envService: s}
}

func (h *DomainHandler) ListByProject(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	domains, err := h.envService.ListDomainsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, domains)
}

func (h *DomainHandler) Create(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	var d models.DomainConfig
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		WriteError(w, http.StatusBadRequest, "domainName is required")
		return nil
	}
	created, err := h.envService.CreateDomain(r.Context(), &d)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *DomainHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	if err := h.envService.DeleteDomain(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

type ProjectEnvHandler struct {
	envService *services.EnvironmentService
}

func NewProjectEnvHandler(s *services.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
}

func (h *ProjectEnvHandler) GetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	vars, err := h.envService.GetVars(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if vars == nil {
		vars = map[string]string{}
	}
	WriteJSON(w, http.StatusOK, vars)
}

func (h *ProjectEnvHandler) SetVars(c echo.Context) error {
	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return nil
	}
	var vars map[string]string
	if err := c.Bind(&vars); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	for k, v := range vars {
		if err := h.envService.SetVar(r.Context(), projectID, k, v); err != nil {
			WriteError(w, http.StatusInternalServerError, err.Error())
			return nil
		}
	}
	WriteJSON(w, http.StatusOK, vars)
}
