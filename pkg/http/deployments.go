package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"vessl.dev/vessl/internal/models"
)

// TriggerDeployment triggers a new manual deployment for an app service.
func (c *Client) TriggerDeployment(serviceID string) (*models.Deployment, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/services/%s/deploy", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetDeploymentStatus checks the status of a specific deployment.
func (c *Client) GetDeploymentStatus(deploymentID string) (*models.Deployment, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/deployments/%s", deploymentID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
