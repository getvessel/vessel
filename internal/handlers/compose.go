package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type ComposeHandler struct {
	projectService  *services.ProjectService
	appService      *services.AppService
	databaseService *services.DatabaseService
	envRepo         repositories.EnvironmentRepository
	appRepo         repositories.AppServiceRepository
	composeParser   *services.ComposeParserService
}

func NewComposeHandler(
	ps *services.ProjectService,
	as *services.AppService,
	ds *services.DatabaseService,
	er repositories.EnvironmentRepository,
	ar repositories.AppServiceRepository,
	cp *services.ComposeParserService,
) *ComposeHandler {
	return &ComposeHandler{
		projectService:  ps,
		appService:      as,
		databaseService: ds,
		envRepo:         er,
		appRepo:         ar,
		composeParser:   cp,
	}
}

type ComposeAnalyzeRequest struct {
	ComposeContent string `json:"composeContent"`
	ProjectID      string `json:"projectId"`
}

func (h *ComposeHandler) Analyze(c echo.Context) error {
	var req ComposeAnalyzeRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}

	if req.ComposeContent == "" {
		return utils.Error(c, http.StatusBadRequest, "compose content is required")
	}

	result, err := h.composeParser.Parse([]byte(req.ComposeContent), req.ProjectID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "failed to parse docker-compose: "+err.Error())
	}

	return utils.Success(c, "Compose analyzed", result)
}

type ComposeDeployRequest struct {
	ProjectID string `json:"projectId"`
}

// @Summary Deploy a docker-compose file
// @Description Parses and deploys all services defined in a docker-compose.yml
// @Tags Compose
// @Accept multipart/form-data
// @Produce json
// @Param projectId formData string false "Project ID (optional, uses default if empty)"
// @Param file formData file true "docker-compose.yml file"
// @Success 200 {object} map[string]any
// @Router /compose/deploy [post]
func (h *ComposeHandler) Deploy(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.FormValue("projectId")
	if projectID == "" {
		projectID = c.FormValue("project_id")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "compose file is required")
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read uploaded file")
	}
	defer src.Close()

	tmpDir := filepath.Join(os.TempDir(), "vessl-compose", uuid.New().String())
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to create temp directory")
	}
	defer os.RemoveAll(tmpDir)

	tmpPath := filepath.Join(tmpDir, "docker-compose.yml")
	dst, err := os.Create(tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save compose file")
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return utils.Error(c, http.StatusInternalServerError, "failed to write compose file")
	}
	dst.Close()

	// Parse the compose file
	composeBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read parsed file")
	}

	result, err := h.composeParser.Parse(composeBytes, projectID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "compose deploy failed: "+err.Error())
	}

	// Create databases first
	var createdCount int
	for _, dbReq := range result.Databases {
		db := &models.Database{
			ProjectID:    dbReq.ProjectID,
			Name:         dbReq.Name,
			Engine:       dbReq.Engine,
			Version:      dbReq.Version,
			Port:         dbReq.Port,
			Username:     dbReq.Username,
			Password:     dbReq.Password,
			DatabaseName: dbReq.DatabaseName,
		}
		_, err := h.databaseService.CreateDatabase(c.Request().Context(), db)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to create database "+dbReq.Name+": "+err.Error())
		}
		createdCount++
	}

	// Create app services
	for _, appReq := range result.AppServices {
		app := &models.AppService{
			ProjectID:      appReq.ProjectID,
			Name:           appReq.Name,
			RuntimeMode:    appReq.RuntimeMode,
			DockerfilePath: appReq.DockerfilePath,
			InstallCommand: appReq.InstallCommand,
			BuildCommand:   appReq.BuildCommand,
			StartCommand:   appReq.StartCommand,
			RepositoryURL:  appReq.RepositoryURL,
			ImageRef:       appReq.ImageRef,
		}
		if appReq.BuildEngine != "" {
			app.BuildEngine = models.BuildEngine(appReq.BuildEngine)
		}

		_, err := h.appService.CreateAppService(c.Request().Context(), app)
		if err != nil {
			return utils.Error(c, http.StatusInternalServerError, "failed to create app service "+app.Name+": "+err.Error())
		}
		createdCount++
	}

	return utils.Success(c, "Compose file deployed", map[string]any{
		"count": createdCount,
	})
}
