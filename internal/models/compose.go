package models

type UserComposeFile struct {
	Services map[string]UserComposeService `yaml:"services" json:"services"`
	Networks map[string]struct {
		External bool `yaml:"external" json:"external"`
	} `yaml:"networks" json:"networks"`
}

type UserComposeService struct {
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Build       any               `yaml:"build,omitempty" json:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty" json:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	EnvFile     string            `yaml:"env_file,omitempty" json:"env_file,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Restart     string            `yaml:"restart,omitempty" json:"restart,omitempty"`
}
