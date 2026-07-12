package templates

type ComposeTemplate struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

type ComposeService struct {
	Image       string          `yaml:"image"`
	Environment []string        `yaml:"environment,omitempty"`
	Ports       []string        `yaml:"ports,omitempty"`
	Volumes     []string        `yaml:"volumes,omitempty"`
	Command     []string        `yaml:"command,omitempty"`
	DependsOn   []string        `yaml:"depends_on,omitempty"`
	XVessl     *VesslMetadata `yaml:"x-vessl,omitempty"`
}

type VesslMetadata struct {
	IsDatabase       bool                  `yaml:"is_database,omitempty"`
	ConnectionString string                `yaml:"connection_string,omitempty"`
	Backup           *VesslBackupMetadata `yaml:"backup,omitempty"`
}

type VesslBackupMetadata struct {
	Command       []string `yaml:"command,omitempty"`
	FileExtension string   `yaml:"file_extension,omitempty"`
}
