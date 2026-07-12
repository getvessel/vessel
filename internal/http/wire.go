//go:build wireinject
// +build wireinject

package http

import (
	"context"
	"database/sql"
	"log"

	"github.com/docker/docker/client"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"vessl.dev/vessl/internal/core"
	"vessl.dev/vessl/internal/engine"
	"vessl.dev/vessl/internal/handlers"
	"vessl.dev/vessl/internal/http/middleware"
	"vessl.dev/vessl/internal/mcp"
	"vessl.dev/vessl/internal/notifications"
	"vessl.dev/vessl/internal/proxy"
	"vessl.dev/vessl/internal/repositories"
	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/vault"
)

func ProvideEcho() *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			log.Printf("REQUEST: %s %s | status: %d", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	return e
}

func ProvideEngineAdapter(
	settings repositories.SettingsRepository,
	service repositories.AppServiceRepository,
	env repositories.EnvRepository,
	database repositories.DatabaseRepository,
	storage repositories.StorageRepository,
	project repositories.ProjectRepository,
	job repositories.JobRepository,
	backup repositories.BackupRepository,
	s3 repositories.S3DestinationRepository,
	svVar repositories.ServiceVarRepository,
	serverless repositories.ServerlessRepository,
) *engineAdapter {
	return newEngineAdapter(settings, service, env, database, storage, project, job, backup, s3, svVar, serverless)
}

func ProvideBackupManagerPath() string {
	return ""
}

func ProvideUpdaterService(s repositories.SettingsRepository) *services.UpdaterService {
	u := services.NewUpdaterService(s)
	u.Start(context.Background())
	return u
}

func ProvideCronManager(c *client.Client, ea engine.CronManagerStore) *engine.CronManager {
	cm := engine.NewCronManager(c, ea)
	_ = cm.Start()
	return cm
}

func ProvideBackupManager(c *client.Client, ea engine.BackupManagerStore) *engine.BackupManager {
	bm := engine.NewBackupManager(c, ea, "")
	_ = bm.Start()
	return bm
}

func ProvideDeploymentListeners(d *core.DispatcherService) *core.DeploymentListeners {
	l := core.NewDeploymentListeners(d)
	l.Register()
	return l
}

var RepoSet = wire.NewSet(
	repositories.NewSettingsSQLiteRepository,
	wire.Bind(new(repositories.SettingsRepository), new(*repositories.SettingsSQLiteRepository)),
	repositories.NewUserSQLiteRepository,
	wire.Bind(new(repositories.UserRepository), new(*repositories.UserSQLiteRepository)),
	repositories.NewOAuthSQLiteRepository,
	wire.Bind(new(repositories.OAuthRepository), new(*repositories.OAuthSQLiteRepository)),
	repositories.NewNotificationSQLiteRepository,
	wire.Bind(new(repositories.NotificationRepository), new(*repositories.NotificationSQLiteRepository)),
	repositories.NewAppServiceSQLiteRepository,
	wire.Bind(new(repositories.AppServiceRepository), new(*repositories.AppServiceSQLiteRepository)),
	repositories.NewGitSQLiteRepository,
	wire.Bind(new(repositories.GitRepository), new(*repositories.GitSQLiteRepository)),
	repositories.NewEnvSQLiteRepository,
	wire.Bind(new(repositories.EnvRepository), new(*repositories.EnvSQLiteRepository)),
	repositories.NewEnvironmentSQLiteRepository,
	wire.Bind(new(repositories.EnvironmentRepository), new(*repositories.EnvironmentSQLiteRepository)),
	repositories.NewDomainSQLiteRepository,
	wire.Bind(new(repositories.DomainRepository), new(*repositories.DomainSQLiteRepository)),
	repositories.NewProjectSQLiteRepository,
	wire.Bind(new(repositories.ProjectRepository), new(*repositories.ProjectSQLiteRepository)),
	repositories.NewDatabaseSQLiteRepository,
	wire.Bind(new(repositories.DatabaseRepository), new(*repositories.DatabaseSQLiteRepository)),
	repositories.NewStorageSQLiteRepository,
	wire.Bind(new(repositories.StorageRepository), new(*repositories.StorageSQLiteRepository)),
	repositories.NewJobSQLiteRepository,
	wire.Bind(new(repositories.JobRepository), new(*repositories.JobSQLiteRepository)),
	repositories.NewCanvasSQLiteRepository,
	wire.Bind(new(repositories.CanvasRepository), new(*repositories.CanvasSQLiteRepository)),
	repositories.NewDeploymentSQLiteRepository,
	wire.Bind(new(repositories.DeploymentRepository), new(*repositories.DeploymentSQLiteRepository)),
	repositories.NewBackupSQLiteRepository,
	wire.Bind(new(repositories.BackupRepository), new(*repositories.BackupSQLiteRepository)),
	repositories.NewTeamSQLiteRepository,
	wire.Bind(new(repositories.TeamRepository), new(*repositories.TeamSQLiteRepository)),
	repositories.NewWorkspaceSQLiteRepository,
	wire.Bind(new(repositories.WorkspaceRepository), new(*repositories.WorkspaceSQLiteRepository)),
	repositories.NewProjectSettingsSQLiteRepository,
	wire.Bind(new(repositories.ProjectSettingsRepository), new(*repositories.ProjectSettingsSQLiteRepository)),
	repositories.NewServiceVarSQLiteRepository,
	wire.Bind(new(repositories.ServiceVarRepository), new(*repositories.ServiceVarSQLiteRepository)),
	repositories.NewS3DestinationSQLiteRepository,
	wire.Bind(new(repositories.S3DestinationRepository), new(*repositories.S3DestinationSQLiteRepository)),
	repositories.NewPRPreviewRepository,
	repositories.NewServerlessRepository,
	repositories.NewGitAppSQLiteRepository,
	wire.Bind(new(repositories.GitAppRepository), new(*repositories.GitAppSQLiteRepository)),
	repositories.NewTeamAISettingsSQLiteRepository,
	wire.Bind(new(repositories.TeamAISettingsRepository), new(*repositories.TeamAISettingsSQLiteRepository)),
	repositories.NewTeamEmailSettingsSQLiteRepository,
	wire.Bind(new(repositories.TeamEmailSettingsRepository), new(*repositories.TeamEmailSettingsSQLiteRepository)),
	repositories.NewVercelRepository,
)

var ServiceSet = wire.NewSet(
	services.NewSettingsService,
	services.NewEmailSettingsService,
	notifications.NewMailerService,
	wire.Bind(new(core.Mailer), new(*notifications.MailerService)),
	core.NewDispatcherService,
	wire.Bind(new(services.NotificationDispatcher), new(*core.DispatcherService)),
	ProvideDeploymentListeners,
	services.NewTokenService,
	ProvideUpdaterService,
	services.NewAppService,
	services.NewGitService,
	services.NewUserService,
	services.NewAuthService,
	services.NewOAuthService,
	services.NewEnvironmentService,
	services.NewProjectService,
	services.NewDatabaseService,
	services.NewStorageService,
	services.NewJobService,
	services.NewCanvasService,
	services.NewBackupService,
	services.NewTeamService,
	services.NewWorkspaceService,
	services.NewProjectSettingsService,
	services.NewNotificationService,
	services.NewDeploymentService,
	services.NewPRPreviewService,
	services.NewGitAppsService,
	services.NewAISettingsService,
	services.NewVercelService,
	services.NewServerlessService,
	services.NewServiceLinker,
)

var EngineSet = wire.NewSet(
	ProvideEngineAdapter,
	wire.Bind(new(engine.CronManagerStore), new(*engineAdapter)),
	wire.Bind(new(engine.BackupManagerStore), new(*engineAdapter)),
	wire.Bind(new(engine.DatabaseDeployerStore), new(*engineAdapter)),
	wire.Bind(new(engine.StorageDeployerStore), new(*engineAdapter)),
	ProvideBackupManagerPath,
	ProvideCronManager,
	ProvideBackupManager,
	engine.NewDatabaseDeployer,
	engine.NewStorageDeployer,
)

var HandlerSet = wire.NewSet(
	handlers.NewAppHandler,
	handlers.NewDatabaseHandler,
	handlers.NewStorageHandler,
	handlers.NewJobHandler,
	handlers.NewCanvasHandler,
	handlers.NewTerminalHandler,
	handlers.NewDeploymentHandler,
	handlers.NewServiceVarHandler,
	handlers.NewProjectSettingsHandler,
	handlers.NewBackupHandler,
	handlers.NewTeamHandler,
	handlers.NewWorkspaceHandler,
	handlers.NewSettingsHandler,
	handlers.NewUpdaterHandler,
	handlers.NewUserHandler,
	handlers.NewAuthHandler,
	handlers.NewOAuthHandler,
	handlers.NewGitHandler,
	handlers.NewWebhookHandler,
	handlers.NewProjectHandler,
	handlers.NewEnvironmentHandler,
	handlers.NewDomainHandler,
	handlers.NewProjectEnvHandler,
	handlers.NewNotificationHandler,
	handlers.NewGitAppsHandler,
	handlers.NewAISettingsHandler,
	handlers.NewEmailSettingsHandler,
	handlers.NewAIDiagnosticsHandler,
	handlers.NewVercelHandler,
	handlers.NewServerlessHandler,
)

var CoreSet = wire.NewSet(
	ProvideEcho,
	mcp.NewBridge,
	middleware.NewAuthGuard,
)

func BuildServer(db *sql.DB, v *vault.Vault, deployer *engine.Deployer, traefikManager *proxy.TraefikManager, dockerClient *client.Client) (*Server, error) {
	wire.Build(
		RepoSet,
		ServiceSet,
		EngineSet,
		HandlerSet,
		CoreSet,
		wire.Bind(new(repositories.Vault), new(*vault.Vault)),
		wire.Bind(new(middleware.SettingsProvider), new(*services.SettingsService)),
		wire.Bind(new(middleware.ProjectTokenProvider), new(*services.ProjectSettingsService)),
		wire.Struct(new(Server), "*"),
	)
	return &Server{}, nil
}
