package engine

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

type ContainerManager struct {
	dockerClient *client.Client
	store        ContainerManagerStore
}

func NewContainerManager(dockerClient *client.Client, st ContainerManagerStore) *ContainerManager {
	return &ContainerManager{dockerClient: dockerClient, store: st}
}

type ContainerRunOptions struct {
	Name            string
	ImageTag        string
	ServiceID       string
	Domain          string
	InternalPort    int
	RuntimeMode     models.RuntimeMode
	Envs            []string
	Cmd             []string
	MemoryLimitMB   int
	CPURequest      float64
	HealthCheckPath string
	Volumes         []models.ServiceVolume
	MaintenanceMode bool
	LogDrains       []*models.LogDrain
}

func (c *ContainerManager) CreateAndStart(ctx context.Context, opts ContainerRunOptions) (string, error) {
	containerPort, err := nat.NewPort("tcp", fmt.Sprintf("%d", opts.InternalPort))
	if err != nil {
		return "", fmt.Errorf("invalid port definition: %w", err)
	}

	if opts.MaintenanceMode {
		opts.ImageTag = "nginx:alpine"
		opts.Cmd = []string{"/bin/sh", "-c", "echo '<html><head><title>Under Maintenance</title><style>body{font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;background:#000;color:#fff;text-align:center;} h1{font-size:2rem;margin-bottom:0.5rem;} p{color:#888;}</style></head><body><div><h1>Under Maintenance</h1><p>This service is temporarily down for maintenance.</p><p>Please check back shortly.</p></div></body></html>' > /usr/share/nginx/html/index.html && nginx -g 'daemon off;'"}
		opts.InternalPort = 80
		opts.HealthCheckPath = "/"
		containerPort, _ = nat.NewPort("tcp", "80")
	}

	config := &container.Config{
		Image: opts.ImageTag,
		Env:   opts.Envs,
		Cmd:   opts.Cmd,
	}

	if opts.HealthCheckPath != "" {
		config.Healthcheck = &container.HealthConfig{
			Test:     []string{"CMD-SHELL", fmt.Sprintf("curl -f http://localhost:%d%s || exit 1", opts.InternalPort, opts.HealthCheckPath)},
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
			Retries:  3,
		}
	}

	if opts.RuntimeMode != models.RuntimeModeWorker {
		config.ExposedPorts = nat.PortSet{containerPort: struct{}{}}
		if opts.ServiceID != "" && opts.Domain != "" {
			config.Labels = c.buildTraefikLabels(opts.ServiceID, opts.Domain, opts.InternalPort, opts.HealthCheckPath)
		}
	}

	var binds []string
	if len(opts.Volumes) > 0 {
		for _, v := range opts.Volumes {
			binds = append(binds, fmt.Sprintf("%s:%s", v.HostPath, v.ContainerPath))
		}
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Resources: container.Resources{
			Memory:   utils.MegaBytesToBytes(opts.MemoryLimitMB),
			NanoCPUs: utils.CPURequestToNanoCPUs(opts.CPURequest),
		},
		NetworkMode: container.NetworkMode(utils.GetRuntimeNetwork()),
		DNS:         c.getCustomDNSResolvers(),
		Binds:       binds,
	}

	_ = c.StopAndRemove(ctx, opts.Name)
	resp, err := c.dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, opts.Name)
	if err != nil {
		return "", fmt.Errorf("docker container create failed: %w", err)
	}
	if err := c.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("docker container start failed: %w", err)
	}

	if len(opts.LogDrains) > 0 {
		StartLogDrains(context.Background(), c.dockerClient, resp.ID, opts.Name, opts.LogDrains)
	}

	return resp.ID, nil
}

func (c *ContainerManager) buildTraefikLabels(serviceID, domain string, internalPort int, healthCheckPath string) map[string]string {
	labels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s.rule", serviceID):                      fmt.Sprintf("Host(`%s`)", domain),
		fmt.Sprintf("traefik.http.routers.%s.tls", serviceID):                       "true",
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", serviceID):          "letsencrypt",
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", serviceID): fmt.Sprintf("%d", internalPort),
	}
	if healthCheckPath != "" {
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.path", serviceID)] = healthCheckPath
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.interval", serviceID)] = "5s"
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.healthcheck.timeout", serviceID)] = "2s"
	}
	return labels
}

func (c *ContainerManager) getCustomDNSResolvers() []string {
	if c.store == nil {
		return nil
	}
	settings, _ := c.store.GetServerSettings()
	if settings == nil || strings.TrimSpace(settings.CustomDNSResolvers) == "" {
		return nil
	}

	parts := strings.Split(settings.CustomDNSResolvers, ",")
	var dnsList []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			dnsList = append(dnsList, p)
		}
	}
	return dnsList
}

func (c *ContainerManager) StopAndRemove(ctx context.Context, containerIDOrName string) error {
	stopTimeout := 10
	_ = c.dockerClient.ContainerStop(ctx, containerIDOrName, container.StopOptions{Timeout: &stopTimeout})
	err := c.dockerClient.ContainerRemove(ctx, containerIDOrName, container.RemoveOptions{Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *ContainerManager) Inspect(ctx context.Context, containerIDOrName string) (types.ContainerJSON, error) {
	return c.dockerClient.ContainerInspect(ctx, containerIDOrName)
}

func (c *ContainerManager) StreamLogs(ctx context.Context, containerIDOrName string, out io.Writer) error {
	logsReader, err := c.dockerClient.ContainerLogs(ctx, containerIDOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	})
	if err != nil {
		return fmt.Errorf("failed to open container logs stream: %w", err)
	}
	defer logsReader.Close()
	_, err = io.Copy(out, logsReader)
	return err
}

func (c *ContainerManager) CleanupOrphanedContainers(ctx context.Context, prefix string, excludeContainerNames []string) error {
	containers, err := c.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	excludeMap := make(map[string]bool)
	for _, n := range excludeContainerNames {
		excludeMap[n] = true
		excludeMap["/"+n] = true
	}

	for _, ctn := range containers {
		for _, name := range ctn.Names {
			if strings.HasPrefix(name, "/"+prefix+"-") {
				if !excludeMap[ctn.ID] && !excludeMap[name] {
					_ = c.StopAndRemove(ctx, ctn.ID)
				}
				break
			}
		}
	}
	return nil
}
