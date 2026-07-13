package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"vessl.dev/vessl/internal/cloud/services"
	"vessl.dev/vessl/internal/utils"
)

// AuthHandler wires HTTP routes to the AuthService.
type AuthHandler struct {
	svc *services.AuthService
}

// NewAuthHandler creates an AuthHandler with the given AuthService.
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles POST /api/cloud/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	token, err := h.svc.Register(c.Request().Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			return utils.Error(c, http.StatusConflict, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Created(c, "User registered", map[string]string{"token": token})
}

// Login handles POST /api/cloud/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	token, err := h.svc.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return utils.Error(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, services.ErrEmailNotVerified) {
			return utils.Error(c, http.StatusForbidden, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Login successful", map[string]string{"token": token})
}

// ForgotPassword handles POST /api/cloud/auth/forgot-password
// Always returns 200 — never reveals whether the email exists.
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	_ = h.svc.ForgotPassword(c.Request().Context(), req.Email)
	return utils.Success(c, "If that email exists, a reset code has been sent", nil)
}

// ResetPassword handles POST /api/cloud/auth/reset-password
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	if err := h.svc.ResetPassword(c.Request().Context(), req.Email, req.OTP, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidOTP) {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Password updated successfully", nil)
}

// VerifyEmail handles GET /api/cloud/auth/verify-email?token=
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return utils.Error(c, http.StatusBadRequest, "missing token")
	}

	if err := h.svc.VerifyEmail(c.Request().Context(), token); err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Email verified successfully", nil)
}
