package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *AuthHandler) Register(c echo.Context) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	u, token, err := h.authService.Register(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	SetAuthCookie(w, token)
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	u, token, err := h.authService.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return nil
	}
	SetAuthCookie(w, token)
	WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	ClearAuthCookie(w)
	WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}
