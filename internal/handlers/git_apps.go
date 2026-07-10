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

func (h *GitAppsHandler) DeleteGithubApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteGithubApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- GitLab Apps ----

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

func (h *GitAppsHandler) DeleteGitlabApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteGitlabApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- Bitbucket Apps ----

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

func (h *GitAppsHandler) DeleteBitbucketApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.gitAppsService.DeleteBitbucketApp(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
