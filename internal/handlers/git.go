package handlers

import (
	"github.com/labstack/echo/v4"

	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type GitHandler struct {
	gitService *services.GitService
}

func NewGitHandler(s *services.GitService) *GitHandler {
	return &GitHandler{gitService: s}
}

func (h *GitHandler) Connect(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	var req models.GitConnectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	gp, err := h.gitService.ConnectProvider(r.Context(), userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, gp)
}

func (h *GitHandler) Status(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	status, err := h.gitService.GetConnectedProviders(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, status)
}

func (h *GitHandler) Disconnect(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	provider := c.Param("provider")
	if provider == "" {
		WriteError(w, http.StatusBadRequest, "missing provider parameter")
		return nil
	}
	if err := h.gitService.DisconnectProvider(r.Context(), userID, provider); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}

func (h *GitHandler) ListRepos(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return nil
	}
	provider := c.QueryParam("provider")
	if provider == "" {
		WriteError(w, http.StatusBadRequest, "missing provider query parameter")
		return nil
	}
	repos, err := h.gitService.ListRepositories(r.Context(), userID, provider)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, repos)
}
