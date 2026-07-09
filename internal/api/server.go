package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/auth"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/oauth"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/store"
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
	gitService             *services.GitService
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
	notifierService        *notifier.NotifierService
	notificationHandler    *NotificationHandler
	updaterService         *updater.UpdaterService
}

// NewServer initializes a Server wired to the database store, container orchestrator, reverse proxy, and Docker client.
func NewServer(s *store.Store, deployer *orchestrator.Deployer, proxyManager *proxy.ProxyManager, dockerClient *client.Client) *Server {
	cronMgr := orchestrator.NewCronManager(dockerClient, s)
	_ = cronMgr.Start()

	backupMgr := orchestrator.NewBackupManager(dockerClient, s, "")
	_ = backupMgr.Start()

	tokenService := services.NewTokenService()
	notifierService := notifier.NewNotifierService(s)

	// Domain repositories
	settingsRepo := settings.NewSQLiteRepository(s.DB())
	userRepo := user.NewSQLiteRepository(s.DB())
	oauthRepo := oauth.NewSQLiteRepository(s.DB())

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
		gitService:             services.NewGitService(s),
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
		notifierService:        notifierService,
		notificationHandler:    NewNotificationHandler(s, notifierService),
		updaterService:         updaterService,
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
