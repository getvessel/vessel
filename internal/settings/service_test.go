package settings_test

import (
	"context"
	"testing"

	"vessel.dev/vessel/internal/settings"
)

type mockRepo struct {
	cfg      *settings.ServerSettings
	projects []map[string]any
}

func (m *mockRepo) GetServerSettings(ctx context.Context) (*settings.ServerSettings, error) {
	return m.cfg, nil
}

func (m *mockRepo) UpdateServerSettings(ctx context.Context, cfg *settings.ServerSettings) error {
	m.cfg = cfg
	return nil
}

func (m *mockRepo) ListProjects(ctx context.Context) ([]map[string]any, error) {
	return m.projects, nil
}

func TestSettingsService_UpdateAndGetSettings(t *testing.T) {
	repo := &mockRepo{
		cfg: &settings.ServerSettings{
			ID:              "global",
			CaddyWildcardIP: "1.2.3.4",
		},
	}
	svc := settings.NewService(repo)

	ctx := context.Background()
	cfg, err := svc.GetSettings(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CaddyWildcardIP != "1.2.3.4" {
		t.Errorf("expected 1.2.3.4, got %s", cfg.CaddyWildcardIP)
	}

	updateCfg := &settings.ServerSettings{
		ID:                  "global",
		CaddyWildcardIP:     "  10.0.0.1  ",
		CustomDNSResolvers:  "  1.1.1.1,8.8.8.8  ",
		RegistrationEnabled: false,
	}
	if err := svc.UpdateSettings(ctx, updateCfg); err != nil {
		t.Fatalf("unexpected error updating: %v", err)
	}

	if repo.cfg.CaddyWildcardIP != "10.0.0.1" {
		t.Errorf("expected trimmed 10.0.0.1, got %s", repo.cfg.CaddyWildcardIP)
	}
	if repo.cfg.CustomDNSResolvers != "1.1.1.1,8.8.8.8" {
		t.Errorf("expected trimmed DNS resolvers, got %s", repo.cfg.CustomDNSResolvers)
	}
}

func TestSettingsService_MCPExecution(t *testing.T) {
	repo := &mockRepo{
		cfg: &settings.ServerSettings{
			MCPServerEnabled: true,
		},
		projects: []map[string]any{{"name": "test-project"}},
	}
	svc := settings.NewService(repo)

	ctx := context.Background()
	res, err := svc.ExecuteMCPTool(ctx, "list_projects")
	if err != nil {
		t.Fatalf("unexpected error executing MCP tool: %v", err)
	}
	if len(res) == 0 {
		t.Errorf("expected non-empty MCP response")
	}

	repo.cfg.MCPServerEnabled = false
	_, err = svc.ExecuteMCPTool(ctx, "list_projects")
	if err == nil {
		t.Errorf("expected error when MCP server is disabled")
	}
}
