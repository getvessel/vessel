package project_settings

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo          Repository
	users         UserProvider
	extractUserID func(r *http.Request) string
}

func NewHandler(repo Repository, users UserProvider, extractUserID func(r *http.Request) string) *Handler {
	return &Handler{
		repo:          repo,
		users:         users,
		extractUserID: extractUserID,
	}
}

func (h *Handler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.repo.ListWebhooksByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var req CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	webhook := &Webhook{
		ProjectID:             projectID,
		URL:                   req.URL,
		EventTypes:            req.EventTypes,
		IncludePREnvironments: req.IncludePREnvironments,
	}

	if err := h.repo.CreateWebhook(r.Context(), webhook); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, webhook)
}

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.repo.DeleteWebhook(r.Context(), id, projectID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTokens(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.repo.ListTokensByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tok := &Token{
		ProjectID:     projectID,
		Name:          req.Name,
		EnvironmentID: req.EnvironmentID,
	}

	fullSecret, err := h.repo.CreateToken(r.Context(), tok)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"token":       tok,
		"secretToken": fullSecret,
	})
}

func (h *Handler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.repo.DeleteToken(r.Context(), id, projectID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	list, err := h.repo.ListMembers(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	member := &ProjectMember{
		ProjectID:  projectID,
		Email:      req.Email,
		Permission: req.Permission,
	}

	if h.users != nil && req.Email != "" {
		user, err := h.users.GetUserByEmail(r.Context(), req.Email)
		if err == nil && user != nil {
			member.UserID = user.ID
		}
	}

	if err := h.repo.AddMember(r.Context(), member); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")

	if err := h.repo.RemoveMember(r.Context(), id, projectID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
