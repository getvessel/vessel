package user

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ClaimsExtractor func(ctx http.Request) string

type Handler struct {
	service       *Service
	extractUserID func(r *http.Request) string
}

func NewHandler(service *Service, extractUserID func(r *http.Request) string) *Handler {
	return &Handler{
		service:       service,
		extractUserID: extractUserID,
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	u, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	u.PasswordHash = ""
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.service.UpdateProfile(r.Context(), userID, req.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.PasswordHash = ""
	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) ListPATs(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	list, err := h.service.ListPATs(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if list == nil {
		list = []*PersonalAccessToken{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) CreatePAT(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	var req CreatePATRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pat, rawToken, err := h.service.CreatePAT(r.Context(), userID, req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

func (h *Handler) DeletePAT(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		writeError(w, http.StatusBadRequest, "missing token id")
		return
	}
	id := pathParts[3]

	if err := h.service.DeletePAT(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
