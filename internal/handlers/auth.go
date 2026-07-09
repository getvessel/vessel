package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	u, token, err := h.authService.Register(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	SetAuthCookie(w, token)
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	u, token, err := h.authService.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	SetAuthCookie(w, token)
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ClearAuthCookie(w)
	WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}
