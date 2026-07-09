package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/auth"
	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/git"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/notification"
	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/oauth"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/project"
	"vessel.dev/vessel/internal/project_env"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/updater"
	"vessel.dev/vessel/internal/user"
)

// Server encapsulates HTTP routing, API handler dependencies, and authentication guards for the Vessel control plane.
type Server struct {
	router                 *http.ServeMux
	store                  *store.Store
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
	deploymentHandler      *DeploymentHandler
	serviceVarHandler      *ServiceVarHandler
	projectSettingsHandler *ProjectSettingsHandler
	backupManager          *orchestrator.BackupManager
	backupHandler          *BackupHandler
	teamHandler            *TeamHandler
	workspaceHandler       *WorkspaceHandler
	settingsHandler        *settings.Handler
	updaterHandler         *updater.Handler
	userHandler            *user.Handler
	authHandler            *auth.Handler
	oauthHandler           *oauth.Handler
	gitHandler             *git.Handler
	projectHandler         *project.Handler
	environmentHandler     *environment.Handler
	domainHandler          *domain.Handler
	projectEnvHandler      *project_env.Handler
	notifierService        *notifier.NotifierService
	notificationHandler    *notification.Handler
	updaterService         *updater.UpdaterService
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	cronMgr := orchestrator.NewCronManager(dockerClient, s)
	_ = cronMgr.Start()

	backupMgr := orchestrator.NewBackupManager(dockerClient, s, "")
	_ = backupMgr.Start()

	tokenService := services.NewTokenService()

	// Domain repositories
	settingsRepo := settings.NewSQLiteRepository(s.DB())
	userRepo := user.NewSQLiteRepository(s.DB())
	oauthRepo := oauth.NewSQLiteRepository(s.DB())
	notifRepo := notification.NewSQLiteRepository(s.DB())

	notifierService := notifier.NewNotifierService(notifRepo)

	// Domain services
	settingsService := settings.NewService(settingsRepo)
	userService := user.NewService(userRepo)
	authService := auth.NewService(userRepo, settingsRepo, tokenService)
	oauthService := oauth.NewService(oauthRepo, userRepo, tokenService)

	updaterService := updater.NewUpdaterService(settingsRepo)
	updaterService.Start(context.Background())

	// Claims extractor helpers
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

	// Git domain
	gitRepo := git.NewSQLiteRepository(s.DB(), s.Vault())
	gitService := git.NewService(gitRepo, nil)
	gitService.WithProjectService(&gitProjectAdapter{store: s})
	gitHandler := git.NewHandler(gitService, extractUserID)

	// Project, environment, domain, and project-env domains
	envRepo := environment.NewSQLiteRepository(s.DB())
	envService := environment.NewService(envRepo)
	envHandler := environment.NewHandler(envService)

	domainRepo := domain.NewSQLiteRepository(s.DB())
	domainService := domain.NewService(domainRepo)
	domainHandler := domain.NewHandler(domainService, proxyManager)

	projectEnvRepo := project_env.NewSQLiteRepository(s.DB(), s.Vault())
	projectEnvService := project_env.NewService(projectEnvRepo)
	projectEnvHandler := project_env.NewHandler(projectEnvService)

	projectRepo := project.NewSQLiteRepository(s.DB(), envRepo)
	projectService := project.NewService(projectRepo, &appServiceRepoAdapter{store: s})
	projectHandler := project.NewHandler(projectService, proxyManager, extractUserID)

	srv := &Server{
		router:                 http.NewServeMux(),
		store:                  s,
		deployer:               deployer,
		proxyManager:           proxyManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              middleware.NewAuthGuard(tokenService, s),
		dbDeployer:             orchestrator.NewDatabaseDeployer(dockerClient, s),
		storageDeployer:        orchestrator.NewStorageDeployer(dockerClient, s),
		cronManager:            cronMgr,
		cronService:            services.NewCronService(s, cronMgr),
		serviceLinker:          services.NewServiceLinker(s),
		deploymentHandler:      NewDeploymentHandler(s),
		serviceVarHandler:      NewServiceVarHandler(s),
		projectSettingsHandler: NewProjectSettingsHandler(s),
		backupManager:          backupMgr,
		backupHandler:          NewBackupHandler(s, backupMgr),
		teamHandler:            NewTeamHandler(s),
		workspaceHandler:       NewWorkspaceHandler(s),
		settingsHandler:        settings.NewHandler(settingsService, dockerClient),
		updaterHandler:         updater.NewHandler(updaterService),
		userHandler:            user.NewHandler(userService, extractUserID),
		authHandler:            auth.NewHandler(authService, extractUserID),
		oauthHandler:           oauth.NewHandler(oauthService, extractClaims),
		gitHandler:             gitHandler,
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

// ServeHTTP satisfies the http.Handler interface, routing through the registered mux with CORS middleware.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.CORSMiddleware(s.router).ServeHTTP(w, r)
}

// Handler returns the root HTTP handler wrapped with global CORS and authentication middleware.
func (s *Server) Handler() http.Handler {
	return middleware.CORSMiddleware(s.router)
}

// RequireAuth validates Bearer tokens or query parameters via middleware before invoking the handler.
func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireAuth(next)
}

// RequireRole enforces that the authenticated user possesses the specified role via middleware.
func (s *Server) RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return s.authGuard.RequireRole(requiredRole, next)
}

// GetUserClaimsFromContext retrieves the authenticated user's claims from request context via middleware.
func GetUserClaimsFromContext(ctx context.Context) *user.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}

// gitProjectAdapter bridges the legacy store app-service query to the git.ProjectService interface.
type gitProjectAdapter struct {
	store *store.Store
}

func (a *gitProjectAdapter) ListAppServicesByProject(projectID string) ([]*git.AppService, error) {
	apps, err := a.store.ListAppServicesByProject(projectID)
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

// appServiceRepoAdapter bridges the legacy store app-service creation to the project.AppServiceRepository interface.
type appServiceRepoAdapter struct {
	store *store.Store
}

func (a *appServiceRepoAdapter) CreateAppService(_ context.Context, app *types.AppServiceConfig) error {
	return a.store.CreateAppService(app)
}
