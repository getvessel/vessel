package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type GitAppsHandler struct {
	gitAppsService *services.GitAppsService
}

func NewGitAppsHandler(gs *services.GitAppsService) *GitAppsHandler {
	return &GitAppsHandler{gitAppsService: gs}
}

// ---- GitHub Apps ----

// @Summary ExchangeGithubManifestCode endpoint
// @Description ExchangeGithubManifestCode endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github/manifest-callback [post]
func (h *GitAppsHandler) ExchangeGithubManifestCode(c echo.Context) error {
	var payload struct {
		Code   string `json:"code"`
		TeamID string `json:"teamId"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if payload.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	app, err := h.gitAppsService.ExchangeGithubManifestCode(c.Request().Context(), payload.Code, payload.TeamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, app)
}

// @Summary ListGithubApps endpoint
// @Description ListGithubApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github [get]
func (h *GitAppsHandler) ListGithubApps(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	apps, err := h.gitAppsService.ListGithubApps(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if apps == nil {
		apps = []models.GithubApp{} // Return empty array instead of null
	}
	return c.JSON(http.StatusOK, apps)
}

// @Summary GetGithubApp endpoint
// @Description GetGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/github/{id} [get]
func (h *GitAppsHandler) GetGithubApp(c echo.Context) error {
	id := c.Param("id")
	app, err := h.gitAppsService.GetGithubApp(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if app == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "App not found"})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary SaveGithubApp endpoint
// @Description SaveGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/github [put]
func (h *GitAppsHandler) SaveGithubApp(c echo.Context) error {
	var app models.GithubApp
	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if app.TeamID == "" {
		app.TeamID = "default"
	}
	if err := h.gitAppsService.SaveGithubApp(c.Request().Context(), &app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary DeleteGithubApp endpoint
// @Description DeleteGithubApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/github/{id} [delete]
func (h *GitAppsHandler) DeleteGithubApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteGithubApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- GitLab Apps ----

// @Summary ListGitlabApps endpoint
// @Description ListGitlabApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/gitlab [get]
func (h *GitAppsHandler) ListGitlabApps(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	apps, err := h.gitAppsService.ListGitlabApps(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if apps == nil {
		apps = []models.GitlabApp{}
	}
	return c.JSON(http.StatusOK, apps)
}

// @Summary GetGitlabApp endpoint
// @Description GetGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/gitlab/{id} [get]
func (h *GitAppsHandler) GetGitlabApp(c echo.Context) error {
	id := c.Param("id")
	app, err := h.gitAppsService.GetGitlabApp(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if app == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "App not found"})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary SaveGitlabApp endpoint
// @Description SaveGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/gitlab [put]
func (h *GitAppsHandler) SaveGitlabApp(c echo.Context) error {
	var app models.GitlabApp
	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if app.TeamID == "" {
		app.TeamID = "default"
	}
	if err := h.gitAppsService.SaveGitlabApp(c.Request().Context(), &app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary DeleteGitlabApp endpoint
// @Description DeleteGitlabApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/gitlab/{id} [delete]
func (h *GitAppsHandler) DeleteGitlabApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteGitlabApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- Bitbucket Apps ----

// @Summary ListBitbucketApps endpoint
// @Description ListBitbucketApps endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/bitbucket [get]
func (h *GitAppsHandler) ListBitbucketApps(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	apps, err := h.gitAppsService.ListBitbucketApps(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if apps == nil {
		apps = []models.BitbucketApp{}
	}
	return c.JSON(http.StatusOK, apps)
}

// @Summary GetBitbucketApp endpoint
// @Description GetBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/bitbucket/{id} [get]
func (h *GitAppsHandler) GetBitbucketApp(c echo.Context) error {
	id := c.Param("id")
	app, err := h.gitAppsService.GetBitbucketApp(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if app == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "App not found"})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary SaveBitbucketApp endpoint
// @Description SaveBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Router /api/settings/git_apps/bitbucket [put]
func (h *GitAppsHandler) SaveBitbucketApp(c echo.Context) error {
	var app models.BitbucketApp
	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if app.TeamID == "" {
		app.TeamID = "default"
	}
	if err := h.gitAppsService.SaveBitbucketApp(c.Request().Context(), &app); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, app)
}

// @Summary DeleteBitbucketApp endpoint
// @Description DeleteBitbucketApp endpoint
// @Tags Settings
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /api/settings/git_apps/bitbucket/{id} [delete]
func (h *GitAppsHandler) DeleteBitbucketApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteBitbucketApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
