package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type OnboardingHandler struct {
	userService    *services.UserService
	authService    *services.AuthService
	settingsRepo   *services.SettingsService
	gitAppsService *services.GitAppsService
	backupService  *services.BackupService
}

func NewOnboardingHandler(
	userService *services.UserService,
	authService *services.AuthService,
	settingsRepo *services.SettingsService,
	gitAppsService *services.GitAppsService,
	backupService *services.BackupService,
) *OnboardingHandler {
	return &OnboardingHandler{
		userService:    userService,
		authService:    authService,
		settingsRepo:   settingsRepo,
		gitAppsService: gitAppsService,
		backupService:  backupService,
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
	cwd, _ := os.Getwd()
	return utils.Success(c, "Setup status", map[string]any{
		"setupRequired": count == 0,
		"cwd":           cwd,
	})
}

type setupEnv struct {
	JWTSecret    string `json:"jwtSecret"`
	DataDir      string `json:"dataDir"`
	DashboardURL string `json:"dashboardUrl"`
	Port         int    `json:"port"`
}

type setupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Env setupEnv `json:"env"`

	DefaultWildcardDomain string `json:"defaultWildcardDomain,omitempty"`
}

// @Summary Complete onboarding setup
// @Description Creates the first user and optionally configures initial settings
// @Tags System
// @Accept json
// @Produce json
// @Param request body setupRequest true "Setup details"
// @Success 200 {object} map[string]any
// @Router /system/setup [post]
func (h *OnboardingHandler) Setup(c echo.Context) error {
	ctx := c.Request().Context()

	count, err := h.userService.CountUsers(ctx)
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	if count > 0 {
		return utils.Error(c, 403, "Setup has already been completed")
	}

	var req setupRequest
	if err := c.Bind(&req); err != nil {
		fmt.Printf("Setup Error: Failed to bind request: %v\n", err)
		return utils.Error(c, 400, "invalid request")
	}

	u, token, err := h.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return utils.Error(c, 400, err.Error())
	}

	if req.Env.JWTSecret == "" {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		req.Env.JWTSecret = hex.EncodeToString(b)
	}

	envBytes, err := os.ReadFile(".env.example")
	if err == nil {
		envStr := string(envBytes)
		envStr = strings.ReplaceAll(envStr, "VESSL_JWT_SECRET=change-this-to-a-secure-random-secret-in-prod", "VESSL_JWT_SECRET="+req.Env.JWTSecret)
		if req.Env.DataDir != "" {
			envStr = strings.ReplaceAll(envStr, "VESSL_DATA_DIR=./data", "VESSL_DATA_DIR="+req.Env.DataDir)
		}
		if req.Env.DashboardURL != "" {
			envStr = strings.ReplaceAll(envStr, "VESSL_DASHBOARD_URL=http://localhost:3000", "VESSL_DASHBOARD_URL="+req.Env.DashboardURL)
		}
		if req.Env.Port != 0 {
			envStr = strings.ReplaceAll(envStr, "PORT=8080", fmt.Sprintf("PORT=%d", req.Env.Port))
		}
		_ = os.WriteFile(".env", []byte(envStr), 0644)
	} else {
		envContent := fmt.Sprintf("VESSL_JWT_SECRET=%s\nVESSL_DATA_DIR=%s\nVESSL_DASHBOARD_URL=%s\nPORT=%d\n",
			req.Env.JWTSecret, req.Env.DataDir, req.Env.DashboardURL, req.Env.Port)
		_ = os.WriteFile(".env", []byte(envContent), 0644)
	}

	settings, err := h.settingsRepo.GetSettings(ctx)
	if err == nil && settings != nil {
		updated := false
		if req.DefaultWildcardDomain != "" {
			settings.DefaultWildcardDomain = req.DefaultWildcardDomain
			updated = true
		}
		if updated {
			_ = h.settingsRepo.UpdateSettings(ctx, settings)
		}
	}

	res := map[string]any{
		"user":  u,
		"token": token,
	}

	return utils.Success(c, "Setup completed successfully", res)
}
