package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	var out []models.User
	for _, u := range users {
		u.PasswordHash = ""
		out = append(out, u)
	}
	WriteJSON(w, http.StatusOK, out)
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	u, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || u == nil {
		WriteError(w, http.StatusNotFound, "user profile not found")
		return nil
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	WriteJSON(w, http.StatusOK, &uCopy)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	u, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || u == nil {
		WriteError(w, http.StatusNotFound, "user profile not found")
		return nil
	}
	var payload struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if payload.Email != "" {
		u.Email = payload.Email
	}
	if payload.Role != "" {
		u.Role = payload.Role
	}
	if err := h.userService.UpdateUser(r.Context(), u); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	uCopy := *u
	uCopy.PasswordHash = ""
	WriteJSON(w, http.StatusOK, &uCopy)
}

func (h *UserHandler) CreatePAT(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
		WriteError(w, http.StatusBadRequest, "token name is required")
		return nil
	}
	pat, rawToken, err := h.userService.CreatePAT(r.Context(), userID, payload.Name, nil)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": rawToken,
		"pat":   pat,
	})
}

func (h *UserHandler) ListPATs(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	pats, err := h.userService.ListPATs(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, pats)
}

func (h *UserHandler) DeletePAT(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	tokenID := c.Param("id")
	if tokenID == "" {
		tokenID = strings.TrimPrefix(r.URL.Path, "/api/auth/pat/")
	}
	if tokenID == "" || tokenID == r.URL.Path {
		WriteError(w, http.StatusBadRequest, "invalid personal access token id")
		return nil
	}
	if err := h.userService.DeletePAT(r.Context(), tokenID, userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
