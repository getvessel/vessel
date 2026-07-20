package engine

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"vessl.dev/vessl/internal/models"
)

func StartLogDrains(ctx context.Context, dockerClient *client.Client, containerID, serviceName string, drains []*models.LogDrain) {
	if len(drains) == 0 {
		return
	}

	go func() {
		time.Sleep(2 * time.Second)

		opts := container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Tail:       "0",
		}

		reader, err := dockerClient.ContainerLogs(context.Background(), containerID, opts)
		if err != nil {
			fmt.Printf("Failed to start log stream for %s: %v\n", containerID, err)
			return
		}
		defer reader.Close()

		hdr := make([]byte, 8)
		for {
			_, err := io.ReadFull(reader, hdr)
			if err != nil {
				break
			}
			count := binary.BigEndian.Uint32(hdr[4:])
			if count > 1024*1024 { // max 1MB per line
				break
			}
			dat := make([]byte, count)
			_, err = io.ReadFull(reader, dat)
			if err != nil {
				break
			}

			logLine := strings.TrimSpace(string(dat))
			if logLine != "" {
				for _, drain := range drains {
					go sendToDrain(drain, serviceName, logLine)
				}
			}
		}
	}()
}

func sendToDrain(drain *models.LogDrain, serviceName, logLine string) {
	payload := map[string]interface{}{
		"service": serviceName,
		"message": logLine,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	switch drain.DrainType {
	case models.LogDrainTypeAxiom:
		sendHTTP(drain.EndpointURL, drain.AuthToken, []interface{}{payload})
	case models.LogDrainTypeNewRelic:
		sendHTTP(drain.EndpointURL, drain.AuthToken, payload)
	case models.LogDrainTypeDatadog:
		ddPayload := map[string]interface{}{
			"ddsource": "vessl",
			"service":  serviceName,
			"message":  logLine,
		}
		sendHTTP(drain.EndpointURL, drain.AuthToken, ddPayload)
	case models.LogDrainTypeWebhook:
		sendHTTP(drain.EndpointURL, drain.AuthToken, payload)
	}
}

func sendHTTP(url, token string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
