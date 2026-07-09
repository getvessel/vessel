package workspace

import (
	"encoding/json"
	"net/http"
)

type ClaimsExtractor func(r *http.Request) (userID, email, role string)

type Handler struct {
	repo          Repository
	extractClaims ClaimsExtractor
}

func NewHandler(repo Repository, extractClaims ClaimsExtractor) *Handler {
	return &Handler{
		repo:          repo,
		extractClaims: extractClaims,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) ListTrustedDomains(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	list, err := h.repo.ListTrustedDomains(r.Context(), teamID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateTrustedDomain(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	var req CreateTrustedDomainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	item := &TrustedDomain{
		TeamID: teamID,
		Domain: req.Domain,
		Role:   req.Role,
	}
	if err := h.repo.CreateTrustedDomain(r.Context(), item); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) DeleteTrustedDomain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	if err := h.repo.DeleteTrustedDomain(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListSSHKeys(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	list, err := h.repo.ListSSHKeys(r.Context(), teamID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateSSHKey(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	var req CreateSSHKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	item := &SSHKey{
		TeamID:    teamID,
		Name:      req.Name,
		PublicKey: req.PublicKey,
	}
	if err := h.repo.CreateSSHKey(r.Context(), item); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	if err := h.repo.DeleteSSHKey(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	teamID := r.PathValue("teamId")
	if teamID == "" {
		writeError(w, http.StatusBadRequest, "missing teamId parameter")
		return
	}

	list, err := h.repo.ListAuditLogs(r.Context(), teamID, 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	ownerID := "default"
	if userID != "" {
		ownerID = userID
	}

	list, err := h.repo.List(r.Context(), ownerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	ownerID := "default"
	if userID != "" {
		ownerID = userID
	}

	var req CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	ws := &Workspace{
		Name:            req.Name,
		AvatarURL:       req.AvatarURL,
		PreferredRegion: req.PreferredRegion,
		OwnerID:         ownerID,
	}
	if err := h.repo.Create(r.Context(), ws); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ws)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	ws, err := h.repo.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}
	writeJSON(w, http.StatusOK, ws)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	ownerID := "default"
	if userID != "" {
		ownerID = userID
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	existing, err := h.repo.Get(r.Context(), id)
	if err != nil || existing == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}
	if existing.OwnerID != ownerID && ownerID != "default" {
		writeError(w, http.StatusForbidden, "permission denied")
		return
	}

	var req UpdateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.AvatarURL != "" {
		existing.AvatarURL = req.AvatarURL
	}
	if req.PreferredRegion != "" {
		existing.PreferredRegion = req.PreferredRegion
	}

	if err := h.repo.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _, _ := h.extractClaims(r)
	ownerID := "default"
	if userID != "" {
		ownerID = userID
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	if err := h.repo.Delete(r.Context(), id, ownerID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
