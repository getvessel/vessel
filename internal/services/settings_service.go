package services

import (
	"context"
	"errors"
	"fmt"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type SettingsService struct {
	settingsRepo     repositories.SettingsRepository
	notificationRepo repositories.NotificationRepository
}

func NewSettingsService(sr repositories.SettingsRepository, nr repositories.NotificationRepository) *SettingsService {
	return &SettingsService{
		settingsRepo:     sr,
		notificationRepo: nr,
	}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*models.ServerSettings, error) {
	return s.settingsRepo.GetServerSettings(ctx)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, cfg *models.ServerSettings) error {
	if cfg == nil {
		return errors.New("server settings cannot be nil")
	}
	return s.settingsRepo.UpdateServerSettings(ctx, cfg)
}

func (s *SettingsService) GetNotificationIntegration(ctx context.Context) (*models.NotificationIntegration, error) {
	return s.notificationRepo.GetIntegration(ctx)
}

func (s *SettingsService) SaveNotificationIntegration(ctx context.Context, n *models.NotificationIntegration) error {
	if n == nil {
		return errors.New("notification integration cannot be nil")
	}
	return s.notificationRepo.SaveIntegration(ctx, n)
}

func (s *SettingsService) GetProjectNotificationPref(ctx context.Context, projectID string) (*models.ProjectNotificationPref, error) {
	if projectID == "" {
		return nil, errors.New("project id is required")
	}
	return s.notificationRepo.GetProjectPref(ctx, projectID)
}

func (s *SettingsService) SaveProjectNotificationPref(ctx context.Context, pref *models.ProjectNotificationPref) error {
	if pref == nil || pref.ProjectID == "" {
		return errors.New("valid project notification preference is required")
	}
	return s.notificationRepo.SaveProjectPref(ctx, pref)
}

func (s *SettingsService) CheckMCPEnabled(ctx context.Context) error {
	settings, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return err
	}
	if settings != nil && !settings.MCPServerEnabled {
		return errors.New("MCP server endpoint is currently disabled by the administrator")
	}
	return nil
}

func (s *SettingsService) ExecuteMCPTool(ctx context.Context, toolName string) ([]map[string]any, error) {
	if err := s.CheckMCPEnabled(ctx); err != nil {
		return nil, err
	}
	switch toolName {
	case "list_projects":
		projects, err := s.settingsRepo.ListProjects(ctx)
		if err != nil {
			return nil, err
		}
		return []map[string]any{
			{"type": "text", "text": fmt.Sprintf("Found %d projects: %+v", len(projects), projects)},
		}, nil
	case "get_system_status":
		return []map[string]any{
			{"type": "text", "text": "Vessel system is healthy and operational."},
		}, nil
	default:
		return nil, fmt.Errorf("Method/Tool not found: %s", toolName)
	}
}
