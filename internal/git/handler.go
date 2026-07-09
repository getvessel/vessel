package git

import (
	"encoding/json"
	"net/http"
)

// Handler serves the HTTP endpoints for git provider management.
type Handler struct {
	service       *Service
	extractUserID func(r *http.Request) string
}

// NewHandler creates a new git Handler.
func NewHandler(service *Service, extractUserID func(r *http.Request) string) *Handler {
	return &Handler{
		service:       service,
		extractUserID: extractUserID,
	}
}

// Connect handles POST /api/git/connect — saves a user's Git provider credentials.
func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req GitConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	gp, err := h.service.SaveProvider(r.Context(), userID, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, gp)
}

// Status handles GET /api/git/status — returns connection status for GitHub and GitLab.
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := h.service.GetConnectedProviders(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// Disconnect handles DELETE /api/git/connect/{provider} — removes a Git provider connection.
func (h *Handler) Disconnect(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider := r.PathValue("provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "missing provider parameter")
		return
	}

	if err := h.service.DisconnectProvider(r.Context(), userID, provider); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}

// ListRepos handles GET /api/git/repos — lists repositories for a given provider.
func (h *Handler) ListRepos(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		writeError(w, http.StatusBadRequest, "missing provider query parameter (e.g. ?provider=github)")
		return
	}

	repos, err := h.service.ListRepositories(r.Context(), userID, provider)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, repos)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
