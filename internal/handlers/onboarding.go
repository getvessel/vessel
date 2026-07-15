package handlers

import (
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type OnboardingHandler struct {
	userService  *services.UserService
	authService  *services.AuthService
	settingsRepo *services.SettingsService
}

func NewOnboardingHandler(userService *services.UserService, authService *services.AuthService, settingsRepo *services.SettingsService) *OnboardingHandler {
	return &OnboardingHandler{
		userService:  userService,
		authService:  authService,
		settingsRepo: settingsRepo,
	}
}

// @Summary Check if onboarding is required
// @Description Returns true if no users exist in the system, indicating setup is needed
// @Tags System
// @Produce json
// @Success 200 {object} map[string]any
// @Router /system/setup-status [get]
func (h *OnboardingHandler) SetupStatus(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	return utils.Success(c, "Setup status", map[string]bool{
		"setupRequired": count == 0,
	})
}

// RegisterRequest defines the expected payload for setup
type SetupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Complete onboarding setup
// @Description Creates the first user and optionally configures initial settings
// @Tags System
// @Accept json
// @Produce json
// @Param request body SetupRequest true "Setup details"
// @Success 200 {object} map[string]any
// @Router /system/setup [post]
func (h *OnboardingHandler) Setup(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	if count > 0 {
		return utils.Error(c, 403, "Setup has already been completed")
	}

	var req SetupRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, 400, "invalid request")
	}

	u, token, err := h.authService.Register(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return utils.Error(c, 400, err.Error())
	}

	res := map[string]any{
		"user":  u,
		"token": token,
	}

	return utils.Success(c, "Setup completed successfully", res)
}
