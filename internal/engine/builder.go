package engine

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/client"
	"golang.org/x/sync/semaphore"

	"vessl.dev/vessl/internal/models"
)

type BuildStrategy string

const (
	StrategyDockerfile BuildStrategy = "dockerfile"
	StrategyRailpack   BuildStrategy = "railpack"
	StrategyNixpacks   BuildStrategy = "nixpacks"
	StrategyBuildpacks BuildStrategy = "buildpacks"
	StrategyServerless BuildStrategy = "serverless"
)

type BuildOptions struct {
	ProjectID      string                `json:"projectId"`
	ServiceID      string                `json:"serviceId,omitempty"`
	SourceDir      string                `json:"sourceDir"`
	DockerfilePath string                `json:"dockerfilePath,omitempty"`
	LogWriter      io.Writer             `json:"-"`
	ProjectConfig  *models.ProjectConfig `json:"projectConfig,omitempty"`
	AppConfig      *models.AppService    `json:"appConfig,omitempty"`
	EnvVars        map[string]string     `json:"envVars,omitempty"`
}

type Builder interface {
	Build(ctx context.Context, opts BuildOptions) (string, error)
}

type EngineBuilder struct {
	dockerClient      *client.Client
	dockerfileBuilder *DockerfileBuilder
	railpackBuilder   *RailpackBuilder
	sem               *semaphore.Weighted
}

func defaultConcurrentBuilds() int64 {
	if s := os.Getenv("VESSL_MAX_CONCURRENT_BUILDS"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v > 0 {
			return v
		}
	}
	return 2
}

func NewBuilder(dockerClient *client.Client) *EngineBuilder {
	return &EngineBuilder{
		dockerClient:      dockerClient,
		dockerfileBuilder: NewDockerfileBuilder(dockerClient),
		railpackBuilder:   NewRailpackBuilder(dockerClient),
		sem:               semaphore.NewWeighted(defaultConcurrentBuilds()),
	}
}

func (b *EngineBuilder) Build(ctx context.Context, opts BuildOptions) (string, error) {
	if err := b.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("failed to acquire build slot: %w", err)
	}
	defer b.sem.Release(1)

	strategy := b.DetectStrategy(opts.SourceDir, opts.DockerfilePath, opts.AppConfig)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🚀 [Builder] Detected build strategy: %s\n", strategy)
	}
	switch strategy {
	case StrategyDockerfile:
		imageTag, err := b.dockerfileBuilder.Build(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("dockerfile build failed: %w", err)
		}
		return imageTag, nil
	case StrategyRailpack, StrategyNixpacks, StrategyBuildpacks:
		imageTag, err := b.railpackBuilder.Build(ctx, opts, string(strategy))
		if err != nil {
			return "", fmt.Errorf("%s build failed: %w", strategy, err)
		}
		return imageTag, nil
	default:
		return "", fmt.Errorf("unsupported build strategy: %s", strategy)
	}
}

func (b *EngineBuilder) DetectStrategy(sourceDir, dockerfilePath string, app *models.AppService) BuildStrategy {
	if app != nil && app.BuildEngine == string(StrategyServerless) {
		return StrategyServerless
	}

	if dockerfilePath != "" {
		if _, err := os.Stat(filepath.Join(sourceDir, dockerfilePath)); err == nil {
			return StrategyDockerfile
		}
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Dockerfile")); err == nil {
		return StrategyDockerfile
	}

	if app != nil && app.BuildEngine != "" && app.BuildEngine != "auto" {
		return BuildStrategy(app.BuildEngine)
	}

	return StrategyRailpack
}
