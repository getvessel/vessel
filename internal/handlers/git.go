package handlers

import (
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

func (h *GitHandler) Connect(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req models.GitConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	gp, err := h.gitService.ConnectProvider(r.Context(), userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, gp)
}

func (h *GitHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	status, err := h.gitService.GetConnectedProviders(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, status)
}

func (h *GitHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	provider := r.PathValue("provider")
	if provider == "" {
		WriteError(w, http.StatusBadRequest, "missing provider parameter")
		return
	}
	if err := h.gitService.DisconnectProvider(r.Context(), userID, provider); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}

func (h *GitHandler) ListRepos(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		WriteError(w, http.StatusBadRequest, "missing provider query parameter")
		return
	}
	repos, err := h.gitService.ListRepositories(r.Context(), userID, provider)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, repos)
}
