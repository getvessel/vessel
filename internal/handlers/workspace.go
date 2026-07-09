package handlers

import (
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

func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	wsList, err := h.workspaceService.ListWorkspaces(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, wsList)
}

func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ws, err := h.workspaceService.CreateWorkspace(r.Context(), payload.Name, userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	ws, err := h.workspaceService.GetWorkspace(r.Context(), id)
	if err != nil || ws == nil {
		WriteError(w, http.StatusNotFound, "workspace not found")
		return
	}
	WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	var ws models.Workspace
	if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ws.ID = id
	if err := h.workspaceService.UpdateWorkspace(r.Context(), &ws); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := ExtractUserID(r)
	if id == "" || userID == "" {
		WriteError(w, http.StatusBadRequest, "missing parameters or unauthorized")
		return
	}
	if err := h.workspaceService.DeleteWorkspace(r.Context(), id, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListTrustedDomains(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}
	domains, err := h.workspaceService.ListTrustedDomains(r.Context(), teamID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, domains)
}

func (h *WorkspaceHandler) CreateTrustedDomain(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}
	var payload struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	td, err := h.workspaceService.AddTrustedDomain(r.Context(), teamID, payload.Domain)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, td)
}

func (h *WorkspaceHandler) DeleteTrustedDomain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	if err := h.workspaceService.DeleteTrustedDomain(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListSSHKeys(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}
	keys, err := h.workspaceService.ListSSHKeys(r.Context(), teamID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, keys)
}

func (h *WorkspaceHandler) CreateSSHKey(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}
	var payload struct {
		Name      string `json:"name"`
		PublicKey string `json:"publicKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	key, err := h.workspaceService.AddSSHKey(r.Context(), teamID, payload.Name, payload.PublicKey)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, key)
}

func (h *WorkspaceHandler) DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	if err := h.workspaceService.DeleteSSHKey(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		WriteError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}
	logs, err := h.workspaceService.ListAuditLogs(r.Context(), teamID, 100)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, logs)
}
