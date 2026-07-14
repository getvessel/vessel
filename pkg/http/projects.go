package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"vessl.dev/vessl/internal/models"
)

// ListProjects returns all projects accessible by the authenticated user.
func (c *Client) ListProjects() ([]*models.ProjectConfig, error) {
	resp, err := c.sendRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list projects (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.ProjectConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(req *models.CreateProjectRequest) (*models.ProjectConfig, error) {
	resp, err := c.sendRequest("POST", "/projects", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create project (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.ProjectConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
