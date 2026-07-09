package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/auth"
	"vessel.dev/vessel/internal/backup"
	"vessel.dev/vessel/internal/canvas"
	"vessel.dev/vessel/internal/database"
	"vessel.dev/vessel/internal/deployment"
	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/env"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/git"
	"vessel.dev/vessel/internal/job"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/notification"
	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/oauth"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/project_settings"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/service_var"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/storage"
	"vessel.dev/vessel/internal/team"
	"vessel.dev/vessel/internal/terminal"
	"vessel.dev/vessel/internal/updater"
	"vessel.dev/vessel/internal/user"
	"vessel.dev/vessel/internal/vault"
	"vessel.dev/vessel/internal/workspace"
)

type Server struct {
	router                 *http.ServeMux
	deployer               *orchestrator.Deployer
	proxyManager           *proxy.ProxyManager
	dockerClient           *client.Client
	tokenService           *services.TokenService
	authGuard              *middleware.AuthGuard
	dbDeployer             *orchestrator.DatabaseDeployer
	storageDeployer        *orchestrator.StorageDeployer
	cronManager            *orchestrator.CronManager
	cronService            *services.CronService
	serviceLinker          *services.ServiceLinker
	gitService             *git.Service
	serviceHandler         *service.Handler
	dbHandler              *database.Handler
	storageHandler         *storage.Handler
	jobHandler             *job.Handler
	canvasHandler          *canvas.Handler
	terminalHandler        *terminal.Handler
	deploymentHandler      *deployment.Handler
	serviceVarHandler      *service_var.Handler
	projectSettingsHandler *project_settings.Handler
	backupHandler          *backup.Handler
	teamHandler            *team.Handler
	workspaceHandler       *workspace.Handler
	settingsHandler        *settings.Handler
	updaterHandler         *updater.Handler
	userHandler            *user.Handler
	authHandler            *auth.Handler
	oauthHandler           *oauth.Handler
	gitHandler             *git.Handler
	webhookHandler         *git.WebhookHandler
	projectHandler         *project.Handler
	environmentHandler     *environment.Handler
	domainHandler          *domain.Handler
	projectEnvHandler      *env.Handler
	notifierService        *notifier.NotifierService
	notificationHandler    *notification.Handler
	updaterService         *updater.UpdaterService
}

func NewServer(db *sql.DB, vault *vault.Vault, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	sa := &storeAdapter{db: db, vault: vault}

	settingsRepo := settings.NewSQLiteRepository(db)
	userRepo := user.NewSQLiteRepository(db)
	oauthRepo := oauth.NewSQLiteRepository(db)
	notifRepo := notification.NewSQLiteRepository(db)

	notifierService := notifier.NewNotifierService(notifRepo)

	settingsService := settings.NewService(settingsRepo)
	userService := user.NewService(userRepo)
	tokenService := services.NewTokenService()
	authService := auth.NewService(userRepo, settingsRepo, tokenService)
	oauthService := oauth.NewService(oauthRepo, userRepo, tokenService)

	updaterService := updater.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())

	extractUserID := func(r *http.Request) string {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID
		}
		return ""
	}
	extractClaims := func(r *http.Request) (userID, email string) {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID, c.Email
		}
		return "", ""
	}
	extractClaims3 := func(r *http.Request) (userID, email, role string) {
		if c := GetUserClaimsFromContext(r.Context()); c != nil {
			return c.UserID, c.Email, c.Role
		}
		return "", "", ""
	}

	serviceRepo := service.NewSQLiteRepository(db)

	gitRepo := git.NewSQLiteRepository(db, vault)
	gitService := git.NewService(gitRepo, nil)
	gitService.WithProjectService(&gitProjectAdapter{svcRepo: serviceRepo})
	gitHandler := git.NewHandler(gitService, extractUserID)

	envRepo := environment.NewSQLiteRepository(db)
	envService := environment.NewService(envRepo)
	envHandler := environment.NewHandler(envService)

	domainRepo := domain.NewSQLiteRepository(db)
	domainService := domain.NewService(domainRepo)
	domainHandler := domain.NewHandler(domainService, proxyManager)

	projectEnvRepo := env.NewSQLiteRepository(db, vault)
	projectEnvService := env.NewService(projectEnvRepo)
	projectEnvHandler := env.NewHandler(projectEnvService)

	projectRepo := project.NewSQLiteRepository(db, envRepo)
	projectService := project.NewService(projectRepo, &appServiceRepoAdapter{svcRepo: serviceRepo})
	projectHandler := project.NewHandler(projectService, proxyManager, extractUserID)

	serviceHandler := service.NewHandler(serviceRepo)

	databaseRepo := database.NewSQLiteRepository(db, vault)
	storageRepo := storage.NewSQLiteRepository(db, vault)
	jobRepo := job.NewSQLiteRepository(db)
	canvasRepo := canvas.NewSQLiteRepository(db, envRepo)
	deploymentRepo := deployment.NewSQLiteRepository(db)
	backupRepo := backup.NewSQLiteRepository(db)
	teamRepo := team.NewSQLiteRepository(db)
	wsRepo := workspace.NewSQLiteRepository(db)
	psRepo := project_settings.NewSQLiteRepository(db)
	svVarRepo := service_var.NewSQLiteRepository(db, &serviceVarSvcAdapter{svcRepo: serviceRepo})

	sa.settingsRepo = settingsRepo
	sa.serviceRepo = serviceRepo
	sa.envRepo = projectEnvRepo
	sa.databaseRepo = databaseRepo
	sa.storageRepo = storageRepo
	sa.projectRepo = projectRepo
	sa.jobRepo = jobRepo
	sa.backupRepo = backupRepo
	sa.deploymentRepo = deploymentRepo
	sa.userRepo = userRepo

	orcCronMgr := orchestrator.NewCronManager(dockerClient, sa)
	_ = orcCronMgr.Start()

	backupMgr := orchestrator.NewBackupManager(dockerClient, sa, "")
	_ = backupMgr.Start()

	dbDeployer := orchestrator.NewDatabaseDeployer(dockerClient, sa)
	storageDeployer := orchestrator.NewStorageDeployer(dockerClient, sa)

	dbHandler := database.NewHandler(databaseRepo, dbDeployer)
	storageHandler := storage.NewHandler(storageRepo, storageDeployer)
	jobHandler := job.NewHandler(jobRepo)
	canvasHandler := canvas.NewHandler(canvasRepo)
	deploymentHandler := deployment.NewHandler(deploymentRepo, &deploymentSvcAdapter{svcRepo: serviceRepo}, &deploymentProjectStoreAdapter{projectRepo: projectRepo}, &deploymentProjectDeployerAdapter{gitService: gitService, deployer: deployer, proxyManager: proxyManager})
	svVarHandler := service_var.NewHandler(svVarRepo, &serviceVarSvcAdapter{svcRepo: serviceRepo})
	terminalHandler := terminal.NewHandler(dockerClient, &tokenValidatorAdapter{inner: tokenService}, &terminalSvcAdapter{svcRepo: serviceRepo})
	backupHandler := backup.NewHandler(backupRepo, backupMgr)
	teamHandler := team.NewHandler(teamRepo, &teamUserProviderAdapter{userRepo: userRepo}, extractClaims3)
	workspaceHandler := workspace.NewHandler(wsRepo, extractClaims3)
	projectSettingsHandler := project_settings.NewHandler(psRepo, &projectSettingsUserProviderAdapter{userRepo: userRepo}, extractUserID)

	srv := &Server{
		router:                 http.NewServeMux(),
		deployer:               deployer,
		proxyManager:           proxyManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              middleware.NewAuthGuard(tokenService, &settingsProviderAdapter{repo: settingsRepo}),
		dbDeployer:             dbDeployer,
		storageDeployer:        storageDeployer,
		cronManager:            orcCronMgr,
		cronService:            services.NewCronService(sa, sa, orcCronMgr),
		serviceLinker:          services.NewServiceLinker(sa, sa),
		serviceHandler:         serviceHandler,
		dbHandler:              dbHandler,
		storageHandler:         storageHandler,
		jobHandler:             jobHandler,
		canvasHandler:          canvasHandler,
		terminalHandler:        terminalHandler,
		deploymentHandler:      deploymentHandler,
		serviceVarHandler:      svVarHandler,
		projectSettingsHandler: projectSettingsHandler,
		backupHandler:          backupHandler,
		teamHandler:            teamHandler,
		workspaceHandler:       workspaceHandler,
		settingsHandler:        settings.NewHandler(settingsService, dockerClient),
		updaterHandler:         updater.NewHandler(updaterService),
		userHandler:            user.NewHandler(userService, extractUserID),
		authHandler:            auth.NewHandler(authService, extractUserID),
		oauthHandler:           oauth.NewHandler(oauthService, extractClaims),
		gitHandler:             gitHandler,
		webhookHandler:         git.NewWebhookHandler(sa, gitService, deployer, proxyManager),
		projectHandler:         projectHandler,
		environmentHandler:     envHandler,
		domainHandler:          domainHandler,
		projectEnvHandler:      projectEnvHandler,
		notifierService:        notifierService,
		notificationHandler: func() *notification.Handler {
			notifService := notification.NewService(notifRepo, notifierService)
			return notification.NewHandler(notifService)
		}(),
		updaterService: updaterService,
	}
	if srv.deployer != nil {
		srv.deployer.EnvProvider = srv.serviceLinker.GetLinkedEnvironmentVariables
	}
	srv.registerRoutes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.CORSMiddleware(s.router).ServeHTTP(w, r)
}

func (s *Server) Handler() http.Handler {
	return middleware.CORSMiddleware(s.router)
}

func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireAuth(next)
}

func (s *Server) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireRole(requiredRole, next)
}

func GetUserClaimsFromContext(ctx context.Context) *user.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}

// ── Legacy adapters ────────────────────────────────────────────────────

type gitProjectAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *gitProjectAdapter) ListAppServicesByProject(projectID string) ([]*git.AppService, error) {
	apps, err := a.svcRepo.ListByProject(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	var result []*git.AppService
	for _, app := range apps {
		result = append(result, &git.AppService{
			ID:            app.ID,
			ProjectID:     app.ProjectID,
			EnvironmentID: app.EnvironmentID,
			Name:          app.Name,
			RepositoryURL: app.RepositoryURL,
			Branch:        app.Branch,
			ContainerID:   app.ContainerID,
		})
	}
	return result, nil
}

type appServiceRepoAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *appServiceRepoAdapter) CreateAppService(ctx context.Context, app *service.AppService) error {
	return a.svcRepo.Create(ctx, app)
}

type deploymentProjectStoreAdapter struct {
	projectRepo *project.SQLiteRepository
}

func (a *deploymentProjectStoreAdapter) GetByID(ctx context.Context, id string) (*deployment.ProjectConfig, error) {
	p, err := a.projectRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return &deployment.ProjectConfig{ID: p.ID, Name: p.Name, Description: p.Description, TeamID: p.TeamID}, nil
}

type deploymentProjectDeployerAdapter struct {
	gitService   *git.Service
	deployer     *orchestrator.Deployer
	proxyManager *proxy.ProxyManager
}

func (a *deploymentProjectDeployerAdapter) CloneOrPullRepository(ctx context.Context, projectID, sourceDir string) error {
	if a.gitService == nil {
		return nil
	}
	return a.gitService.CloneOrPullRepository(ctx, projectID, sourceDir, nil)
}

func (a *deploymentProjectDeployerAdapter) DeployProject(ctx context.Context, cfg *deployment.ProjectConfig, sourceDir string) (string, error) {
	if a.deployer == nil {
		return "", fmt.Errorf("deployer not available")
	}
	p := &project.ProjectConfig{ID: cfg.ID, Name: cfg.Name, Description: cfg.Description, TeamID: cfg.TeamID}
	return a.deployer.Deploy(ctx, p, sourceDir, nil)
}

func (a *deploymentProjectDeployerAdapter) ReloadProxy(ctx context.Context) error {
	if a.proxyManager == nil {
		return nil
	}
	return a.proxyManager.Reload(ctx)
}

type deploymentSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *deploymentSvcAdapter) GetByID(ctx context.Context, id string) (any, error) {
	return a.svcRepo.GetByID(ctx, id)
}

type serviceVarSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *serviceVarSvcAdapter) GetByID(ctx context.Context, id string) (*service_var.ServiceDTO, error) {
	app, err := a.svcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return &service_var.ServiceDTO{
		ID:            app.ID,
		ProjectID:     app.ProjectID,
		EnvironmentID: app.EnvironmentID,
	}, nil
}

type terminalSvcAdapter struct {
	svcRepo *service.SQLiteRepository
}

func (a *terminalSvcAdapter) GetByID(ctx context.Context, id string) (*terminal.AppService, error) {
	app, err := a.svcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return &terminal.AppService{ID: app.ID, ContainerID: app.ContainerID}, nil
}

type teamUserProviderAdapter struct {
	userRepo *user.SQLiteRepository
}

func (a *teamUserProviderAdapter) GetUserByEmail(email string) (*user.User, error) {
	return a.userRepo.GetUserByEmail(context.Background(), email)
}

type projectSettingsUserProviderAdapter struct {
	userRepo *user.SQLiteRepository
}

func (a *projectSettingsUserProviderAdapter) GetUserByEmail(ctx context.Context, email string) (*project_settings.User, error) {
	u, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return &project_settings.User{ID: u.ID, Email: u.Email}, nil
}

type tokenValidatorAdapter struct {
	inner *services.TokenService
}

func (a *tokenValidatorAdapter) ValidateToken(tokenStr string) (*terminal.TokenClaim, error) {
	claims, err := a.inner.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}
	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	return &terminal.TokenClaim{UserID: sub, Email: email}, nil
}

type settingsProviderAdapter struct {
	repo *settings.SQLiteRepository
}

func (a *settingsProviderAdapter) GetServerSettings() (*settings.ServerSettings, error) {
	return a.repo.GetServerSettings(context.Background())
}

// storeAdapter implements all legacy store interfaces consumed by orchestrator, services, and git.
type storeAdapter struct {
	db             *sql.DB
	vault          *vault.Vault
	settingsRepo   *settings.SQLiteRepository
	serviceRepo    *service.SQLiteRepository
	envRepo        *env.SQLiteRepository
	databaseRepo   *database.SQLiteRepository
	storageRepo    *storage.SQLiteRepository
	projectRepo    *project.SQLiteRepository
	jobRepo        *job.SQLiteRepository
	backupRepo     *backup.SQLiteRepository
	deploymentRepo *deployment.SQLiteRepository
	userRepo       *user.SQLiteRepository
}

// ── ContainerManagerStore ─────────────────────────────────────────────

func (a *storeAdapter) GetServerSettings() (*settings.ServerSettings, error) {
	return a.settingsRepo.GetServerSettings(context.Background())
}

// ── DeployerStore ────────────────────────────────────────────────────

func (a *storeAdapter) ListAppServicesByProject(projectID string) ([]*service.AppService, error) {
	return a.serviceRepo.ListByProject(context.Background(), projectID)
}

func (a *storeAdapter) GetEnvVars(projectID string) (map[string]string, error) {
	return a.envRepo.GetVars(context.Background(), projectID)
}

func (a *storeAdapter) ListServiceVariables(serviceID string) ([]*service_var.Variable, error) {
	svVarRepo := service_var.NewSQLiteRepository(a.db, &serviceVarSvcAdapter{svcRepo: a.serviceRepo})
	return svVarRepo.ListByService(context.Background(), serviceID)
}

// ── DatabaseDeployerStore ────────────────────────────────────────────

func (a *storeAdapter) UpdateDatabaseStatus(id string, status string, containerID string) error {
	_, err := a.db.Exec(`UPDATE databases SET status = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, containerID, id)
	return err
}

func (a *storeAdapter) GetDatabase(id string) (*database.Database, error) {
	return a.databaseRepo.GetByID(context.Background(), id)
}

// ── StorageDeployerStore ─────────────────────────────────────────────

func (a *storeAdapter) UpdateStorageStatus(id string, status string, containerID string) error {
	_, err := a.db.Exec(`UPDATE storage SET status = ?, container_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, containerID, id)
	return err
}

func (a *storeAdapter) GetStorage(id string) (*storage.Storage, error) {
	return a.storageRepo.GetByID(context.Background(), id)
}

// ── CronManagerStore ─────────────────────────────────────────────────

func (a *storeAdapter) ListJobs() ([]job.Job, error) {
	return a.jobRepo.ListByProject(context.Background(), "")
}

func (a *storeAdapter) GetJob(id string) (*job.Job, error) {
	return a.jobRepo.GetByID(context.Background(), id)
}

func (a *storeAdapter) GetProject(id string) (*project.ProjectConfig, error) {
	return a.projectRepo.Get(context.Background(), id)
}

func (a *storeAdapter) UpdateJobStatusAndOutput(id string, status string, lastRunAt *time.Time, output string) error {
	return a.jobRepo.UpdateStatus(context.Background(), id, status, lastRunAt, output)
}

// ── BackupManagerStore ───────────────────────────────────────────────

func (a *storeAdapter) ListAllActiveBackupConfigs() ([]*backup.BackupConfig, error) {
	return a.backupRepo.ListAllActiveConfigs(context.Background())
}

func (a *storeAdapter) GetBackupConfig(id string) (*backup.BackupConfig, error) {
	return a.backupRepo.GetConfigByID(context.Background(), id)
}

func (a *storeAdapter) CreateBackupRecord(rec *backup.BackupRecord) error {
	return a.backupRepo.CreateRecord(context.Background(), rec)
}

func (a *storeAdapter) UpdateBackupRecord(id, status, filePath, s3URL, logs string, fileSizeBytes int64, completedAt string) error {
	_, err := a.db.Exec(`UPDATE backup_records SET status = ?, file_path = ?, s3_url = ?, logs = ?, file_size_bytes = ?, completed_at = ? WHERE id = ?`,
		status, filePath, s3URL, logs, fileSizeBytes, completedAt, id)
	return err
}

func (a *storeAdapter) GetS3Destination(id string) (*backup.S3Destination, error) {
	return a.backupRepo.GetS3Destination(context.Background(), id)
}

func (a *storeAdapter) ListBackupRecords(backupConfigID string) ([]*backup.BackupRecord, error) {
	return a.backupRepo.ListRecordsByConfig(context.Background(), backupConfigID)
}

// ── services.JobStore ────────────────────────────────────────────────

func (a *storeAdapter) CreateJob(j *job.Job) error {
	return a.jobRepo.Create(context.Background(), j)
}

func (a *storeAdapter) ListJobsByProject(projectID string) ([]job.Job, error) {
	return a.jobRepo.ListByProject(context.Background(), projectID)
}

func (a *storeAdapter) DeleteJob(id string) error {
	return a.jobRepo.Delete(context.Background(), id)
}

// ── services.ProjectStore ────────────────────────────────────────────

// GetProject is already defined above for CronManagerStore.

// ── services.DatabaseLister ──────────────────────────────────────────

func (a *storeAdapter) ListDatabasesByProject(projectID string) ([]database.Database, error) {
	rows, err := a.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, engine, version, port, username, encrypted_password, database_name, volume_path, COALESCE(container_id, ''), status, COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM databases WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []database.Database
	for rows.Next() {
		var d database.Database
		var encryptedPassword string
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.EnvironmentID, &d.Name, &d.Engine, &d.Version, &d.Port, &d.Username, &encryptedPassword, &d.DatabaseName, &d.VolumePath, &d.ContainerID, &d.Status, &d.InternalDNS, &d.ExternalDNS, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

// ── services.StorageLister ───────────────────────────────────────────

func (a *storeAdapter) ListStorageByProject(projectID string) ([]storage.Storage, error) {
	rows, err := a.db.Query(`SELECT id, COALESCE(project_id, ''), COALESCE(environment_id, ''), name, type, api_port, console_port, access_key, encrypted_secret_key, bucket_name, COALESCE(volume_path, ''), COALESCE(container_id, ''), COALESCE(status, 'stopped'), COALESCE(internal_dns, ''), COALESCE(external_dns, ''), created_at, updated_at FROM storage WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []storage.Storage
	for rows.Next() {
		var s storage.Storage
		var encryptedSecretKey string
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.EnvironmentID, &s.Name, &s.Type, &s.APIPort, &s.ConsolePort, &s.AccessKey, &encryptedSecretKey, &s.BucketName, &s.VolumePath, &s.ContainerID, &s.Status, &s.InternalDNS, &s.ExternalDNS, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// ── git.Store ────────────────────────────────────────────────────────

func (a *storeAdapter) GetAppService(id string) (*service.AppService, error) {
	return a.serviceRepo.GetByID(context.Background(), id)
}

func (a *storeAdapter) CreateDeployment(dep *deployment.Deployment) error {
	return a.deploymentRepo.Create(context.Background(), dep)
}

func (a *storeAdapter) UpdateDeploymentStatus(id, status, buildLogs, containerID string) error {
	return a.deploymentRepo.UpdateStatus(context.Background(), id, status, buildLogs, containerID)
}
