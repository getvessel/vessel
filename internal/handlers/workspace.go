package handlers

import (
	"github.com/labstack/echo/v4"

	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type WorkspaceHandler struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceHandler(s *services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: s}
}

func (h *WorkspaceHandler) List(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	wsList, err := h.workspaceService.ListWorkspaces(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, wsList)
}

func (h *WorkspaceHandler) Create(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	ws, err := h.workspaceService.CreateWorkspace(r.Context(), payload.Name, userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	ws, err := h.workspaceService.GetWorkspace(r.Context(), id)
	if err != nil || ws == nil {
		WriteError(w, http.StatusNotFound, "workspace not found")
		return nil
	}
	WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	var ws models.Workspace
	if err := c.Bind(&ws); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	ws.ID = id
	if err := h.workspaceService.UpdateWorkspace(r.Context(), &ws); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	userID := ExtractUserID(r)
	if id == "" || userID == "" {
		WriteError(w, http.StatusBadRequest, "missing parameters or unauthorized")
		return nil
	}
	if err := h.workspaceService.DeleteWorkspace(r.Context(), id, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListTrustedDomains(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return nil
	}
	domains, err := h.workspaceService.ListTrustedDomains(r.Context(), teamID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, domains)
}

func (h *WorkspaceHandler) CreateTrustedDomain(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return nil
	}
	var payload struct {
		Domain string `json:"domain"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	td, err := h.workspaceService.AddTrustedDomain(r.Context(), teamID, payload.Domain)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, td)
}

func (h *WorkspaceHandler) DeleteTrustedDomain(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	if err := h.workspaceService.DeleteTrustedDomain(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListSSHKeys(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return nil
	}
	keys, err := h.workspaceService.ListSSHKeys(r.Context(), teamID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, keys)
}

func (h *WorkspaceHandler) CreateSSHKey(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return nil
	}
	var payload struct {
		Name      string `json:"name"`
		PublicKey string `json:"publicKey"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	key, err := h.workspaceService.AddSSHKey(r.Context(), teamID, payload.Name, payload.PublicKey)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, key)
}

func (h *WorkspaceHandler) DeleteSSHKey(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return nil
	}
	if err := h.workspaceService.DeleteSSHKey(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListAuditLogs(c echo.Context) error {
	teamID := c.Param("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return nil
	}
	logs, err := h.workspaceService.ListAuditLogs(r.Context(), teamID, 100)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, logs)
}
