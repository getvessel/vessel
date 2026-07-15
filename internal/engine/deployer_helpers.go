package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/utils"
)

func (d *Deployer) getEnvironmentVariables(app *models.AppService, logWriter io.Writer) (map[string]string, error) {
	envVarsMap, err := d.store.GetEnvVars(app.ProjectID)
	if err != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "⚠️ [Deployer] Warning: could not load shared project environment variables: %v\n", err)
	}
	if envVarsMap == nil {
		envVarsMap = make(map[string]string)
	}

	serviceVars, _ := d.store.ListServiceVariables(app.ID)
	for _, sv := range serviceVars {
		envVarsMap[sv.Key] = sv.Value
	}

	if d.EnvProvider != nil {
		if linkedEnvs, err := d.EnvProvider(app.ProjectID); err == nil {
			for k, v := range linkedEnvs {
				if _, exists := envVarsMap[k]; !exists {
					envVarsMap[k] = v
				}
			}
			if logWriter != nil && len(linkedEnvs) > 0 {
				fmt.Fprintf(logWriter, "🔗 [Deployer] Automatically linked %d service connection strings (DATABASE_URL, REDIS_URL, etc.)\n", len(linkedEnvs))
			}
		}
	}

	if d.EnvInterpolator != nil {
		if registry, err := d.EnvInterpolator(app.ProjectID); err == nil && len(registry) > 0 {
			envVarsMap = InterpolateEnvVars(envVarsMap, registry)
			if logWriter != nil {
				fmt.Fprintf(logWriter, "🔀 [Deployer] Interpolated dynamic variable references (${service.VAR_KEY} syntax).\n")
			}
		}
	}

	return envVarsMap, nil
}

func defaultAppPort() int {
	if p := os.Getenv("VESSL_DEFAULT_APP_PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil && port > 0 {
			return port
		}
	}
	return 3000
}

func defaultMemoryMB() int {
	if m := os.Getenv("VESSL_DEFAULT_MEMORY_MB"); m != "" {
		if mem, err := strconv.Atoi(m); err == nil && mem > 0 {
			return mem
		}
	}
	return 512
}

func defaultCPURequest() float64 {
	if c := os.Getenv("VESSL_DEFAULT_CPU"); c != "" {
		if cpu, err := strconv.ParseFloat(c, 64); err == nil && cpu > 0 {
			return cpu
		}
	}
	return 0.5
}

func (d *Deployer) verifyHealthCheck(ctx context.Context, app *models.AppService, containerName string, logWriter io.Writer) error {
	if app.RuntimeMode == models.RuntimeModeWorker {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "✅ [Deployer] Worker mode detected. Skipping HTTP health check.\n")
		}
		return nil
	}
	healthy := d.waitForHealthyContainer(ctx, containerName, app.HealthCheckPath, app.InternalPort)
	if !healthy {
		_ = d.containerManager.StopAndRemove(ctx, containerName)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "❌ [Deployer] Health check failed. Rolling back to previous version.\n")
		}
		return fmt.Errorf("health check failed, deployment aborted")
	}
	return nil
}

func (d *Deployer) waitForHealthyContainer(ctx context.Context, containerName string, healthCheckPath string, internalPort int) bool {
	maxRetries := 30
	if t := os.Getenv("VESSL_DEPLOYMENT_TIMEOUT"); t != "" {
		if v, err := strconv.Atoi(t); err == nil && v > 0 {
			maxRetries = v / 2
		}
	}
	for i := 0; i < maxRetries; i++ {
		time.Sleep(2 * time.Second)
		inspect, err := d.containerManager.Inspect(ctx, containerName)
		if err == nil {
			if !inspect.State.Running {
				if inspect.State.Status == "exited" {
					break
				}
				continue
			}
			if healthCheckPath != "" {
				var containerIP string
				if net, ok := inspect.NetworkSettings.Networks[utils.GetRuntimeNetwork()]; ok {
					containerIP = net.IPAddress
				}
				if containerIP != "" {
					port := internalPort
					if port <= 0 {
						port = 3000
					}
					resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", containerIP, port, healthCheckPath))
					if err == nil {
						resp.Body.Close()
						if resp.StatusCode >= 200 && resp.StatusCode < 400 {
							return true
						}
					}
				}
			} else {
				return true
			}
		}
	}
	return false
}
