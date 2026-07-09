package settings

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSettings(ctx context.Context) (*ServerSettings, error) {
	return s.repo.GetServerSettings(ctx)
}

func (s *Service) UpdateSettings(ctx context.Context, cfg *ServerSettings) error {
	if cfg == nil {
		return errors.New("settings configuration cannot be nil")
	}
	cfg.CustomDNSResolvers = strings.TrimSpace(cfg.CustomDNSResolvers)
	cfg.IPAllowlist = strings.TrimSpace(cfg.IPAllowlist)
	cfg.CaddyWildcardIP = strings.TrimSpace(cfg.CaddyWildcardIP)
	if cfg.CaddyWildcardIP == "" {
		cfg.CaddyWildcardIP = "127.0.0.1"
	}
	return s.repo.UpdateServerSettings(ctx, cfg)
}

func (s *Service) CheckMCPEnabled(ctx context.Context) error {
	settings, err := s.repo.GetServerSettings(ctx)
	if err != nil {
		return err
	}
	if settings != nil && !settings.MCPServerEnabled {
		return errors.New("MCP server endpoint is currently disabled by the administrator")
	}
	return nil
}

func (s *Service) ExecuteMCPTool(ctx context.Context, toolName string) ([]map[string]any, error) {
	if err := s.CheckMCPEnabled(ctx); err != nil {
		return nil, err
	}
	switch toolName {
	case "list_projects":
		projects, err := s.repo.ListProjects(ctx)
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
