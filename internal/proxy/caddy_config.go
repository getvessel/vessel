package proxy

import (
	"os"
	"path/filepath"
)

// CaddyConfig encapsulates global file paths, network settings, and administration endpoints for the Caddy reverse proxy.
type CaddyConfig struct {
	CaddyfilePath    string
	AdminAPIEndpoint string
	DockerContainer  string
	TLSEmail         string
}

// NewCaddyConfig returns a CaddyConfig initialized with standard default paths and administrative endpoints.
func NewCaddyConfig(baseDataDir string, tlsEmail string) *CaddyConfig {
	if baseDataDir == "" {
		baseDataDir = "data"
	}
	caddyDir := filepath.Join(baseDataDir, "caddy")
	_ = os.MkdirAll(caddyDir, 0755)

	if tlsEmail == "" {
		tlsEmail = os.Getenv("VESSEL_TLS_EMAIL")
	}

	return &CaddyConfig{
		CaddyfilePath:    filepath.Join(caddyDir, "Caddyfile"),
		AdminAPIEndpoint: "http://127.0.0.1:2019/load",
		DockerContainer:  "vessel-caddy",
		TLSEmail:         tlsEmail,
	}
}
