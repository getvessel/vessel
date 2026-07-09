package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var out []models.User
	for _, u := range users {
		u.PasswordHash = ""
		out = append(out, u)
	}
	WriteJSON(w, http.StatusOK, out)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	u, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || u == nil {
		WriteError(w, http.StatusNotFound, "user profile not found")
		return
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	WriteJSON(w, http.StatusOK, &uCopy)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	u, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || u == nil {
		WriteError(w, http.StatusNotFound, "user profile not found")
		return
	}
	var payload struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if payload.Email != "" {
		u.Email = payload.Email
	}
	if payload.Role != "" {
		u.Role = payload.Role
	}
	if err := h.userService.UpdateUser(r.Context(), u); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	WriteJSON(w, http.StatusOK, &uCopy)
}

func (h *UserHandler) CreatePAT(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		WriteError(w, http.StatusBadRequest, "token name is required")
		return
	}
	pat, rawToken, err := h.userService.CreatePAT(r.Context(), userID, payload.Name, nil)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

func (h *UserHandler) ListPATs(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	pats, err := h.userService.ListPATs(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, pats)
}

func (h *UserHandler) DeletePAT(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	tokenID := strings.TrimPrefix(r.URL.Path, "/api/auth/pat/")
	if tokenID == "" || tokenID == r.URL.Path {
		WriteError(w, http.StatusBadRequest, "invalid personal access token id")
		return
	}
	if err := h.userService.DeletePAT(r.Context(), tokenID, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
