package types

import "time"

type ContainerHealth struct {
	Status             string  `json:"status"`
	CPUUsagePercentage float64 `json:"cpuUsagePercentage"`
	MemoryUsageBytes   int64   `json:"memoryUsageBytes"`
	MemoryLimitBytes   int64   `json:"memoryLimitBytes"`
	UptimeSeconds      int64   `json:"uptimeSeconds"`
}

type ProjectConfig struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	RepositoryURL  string          `json:"repositoryUrl,omitempty"`
	Branch         string          `json:"branch,omitempty"`
	DockerfilePath string          `json:"dockerfilePath,omitempty"`
	Domain         string          `json:"domain,omitempty"`
	EnvVarsCount   int             `json:"envVarsCount"`
	Health         ContainerHealth `json:"health"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type SystemInfo struct {
	Version         string `json:"version"`
	GoVersion       string `json:"goVersion"`
	DockerVersion   string `json:"dockerVersion"`
	CaddyVersion    string `json:"caddyVersion"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	TotalMemoryMB   int64  `json:"totalMemoryMB"`
	FreeMemoryMB    int64  `json:"freeMemoryMB"`
	CPUCores        int    `json:"cpuCores"`
	UpdateAvailable bool   `json:"updateAvailable"`
	LatestVersion   string `json:"latestVersion,omitempty"`
}
