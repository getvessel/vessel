package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/docker/docker/client"
	"vessel.dev/vessel/internal/api"
	"vessel.dev/vessel/internal/orchestrator"
	"vessel.dev/vessel/internal/proxy"
	"vessel.dev/vessel/internal/store"
)

const vesselVersion = "0.1.0-alpha"

func main() {
	log.Printf("🛰️ Booting Vessel Daemon (`vesseld`) v%s [%s/%s]...", vesselVersion, runtime.GOOS, runtime.GOARCH)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataDir := os.Getenv("VESSEL_DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	dbStore, err := store.NewStore(dataDir)
	if err != nil {
		log.Fatalf("❌ Failed to initialize state store (`%s/vessel.db`): %v", dataDir, err)
	}
	defer dbStore.Close()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("⚠️ Docker daemon connection warning: %v (container deployment features disabled)", err)
	}

	deployer := orchestrator.NewDeployer(dockerClient, dbStore)
	proxyCfg := proxy.NewCaddyConfig(dataDir, os.Getenv("VESSEL_TLS_EMAIL"))
	proxyMgr := proxy.NewProxyManager(proxyCfg, dbStore, dockerClient)

	_ = proxyMgr.Reload(context.Background())

	apiServer := api.NewServer(dbStore, deployer, proxyMgr, dockerClient)

	log.Printf("🚀 Vessel control plane listening on :%s", port)
	if err := http.ListenAndServe(":"+port, apiServer.Handler()); err != nil {
		log.Fatalf("❌ Server crashed: %v", err)
	}
}
