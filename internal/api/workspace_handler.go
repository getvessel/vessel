package api

import (
	"encoding/json"
	"net/http"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

type WorkspaceHandler struct {
	store *store.Store
}

func NewWorkspaceHandler(store *store.Store) *WorkspaceHandler {
	return &WorkspaceHandler{store: store}
}

// ListTrustedDomains lists all trusted SSO/workspace domains for a team.
func (h *WorkspaceHandler) ListTrustedDomains(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceTrustedDomains(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// CreateTrustedDomain adds a trusted SSO/workspace domain.
func (h *WorkspaceHandler) CreateTrustedDomain(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	var item types.WorkspaceTrustedDomain
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.TeamID = teamID

	if err := h.store.CreateWorkspaceTrustedDomain(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// DeleteTrustedDomain removes a trusted domain.
func (h *WorkspaceHandler) DeleteTrustedDomain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteWorkspaceTrustedDomain(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListSSHKeys lists all SSH public keys registered for a workspace team.
func (h *WorkspaceHandler) ListSSHKeys(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceSSHKeys(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// CreateSSHKey registers a new SSH public key for the workspace.
func (h *WorkspaceHandler) CreateSSHKey(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	var item types.WorkspaceSSHKey
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item.TeamID = teamID

	if err := h.store.CreateWorkspaceSSHKey(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// DeleteSSHKey removes an SSH key from the workspace.
func (h *WorkspaceHandler) DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteWorkspaceSSHKey(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListAuditLogs retrieves the audit log history for a workspace team.
func (h *WorkspaceHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		http.Error(w, "missing teamId parameter", http.StatusBadRequest)
		return
	}

	list, err := h.store.ListWorkspaceAuditLogs(teamID, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
