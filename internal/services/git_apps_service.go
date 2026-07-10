package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type GitAppsService struct {
	repo repositories.GitAppRepository
}

func NewGitAppsService(repo repositories.GitAppRepository) *GitAppsService {
	return &GitAppsService{repo: repo}
}

// ---- GitHub Apps ----

type githubManifestConversionResponse struct {
	ID            int    `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
	HTMLURL       string `json:"html_url"`
	Name          string `json:"name"`
}

func (s *GitAppsService) ExchangeGithubManifestCode(ctx context.Context, code string, teamID string) (*models.GithubApp, error) {
	if code == "" {
		return nil, errors.New("conversion code is required")
	}
	if teamID == "" {
		teamID = "default"
	}

	url := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api error: %s", string(body))
	}

	var conversion githubManifestConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&conversion); err != nil {
		return nil, err
	}

	app := &models.GithubApp{
		ID:            uuid.NewString(),
		TeamID:        teamID,
		Name:          conversion.Name,
		AppID:         fmt.Sprintf("%d", conversion.ID),
		ClientID:      conversion.ClientID,
		ClientSecret:  conversion.ClientSecret,
		WebhookSecret: conversion.WebhookSecret,
		PrivateKey:    conversion.PEM,
		IsPublic:      false,
	}

	if err := s.repo.SaveGithubApp(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *GitAppsService) ListGithubApps(ctx context.Context, teamID string) ([]models.GithubApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListGithubApps(ctx, teamID)
}

func (s *GitAppsService) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetGithubApp(ctx, id)
}

func (s *GitAppsService) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveGithubApp(ctx, app)
}

func (s *GitAppsService) DeleteGithubApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteGithubApp(ctx, id)
}

// ---- GitLab Apps ----

func (s *GitAppsService) ListGitlabApps(ctx context.Context, teamID string) ([]models.GitlabApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListGitlabApps(ctx, teamID)
}

func (s *GitAppsService) GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetGitlabApp(ctx, id)
}

func (s *GitAppsService) SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveGitlabApp(ctx, app)
}

func (s *GitAppsService) DeleteGitlabApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteGitlabApp(ctx, id)
}

// ---- Bitbucket Apps ----

func (s *GitAppsService) ListBitbucketApps(ctx context.Context, teamID string) ([]models.BitbucketApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListBitbucketApps(ctx, teamID)
}

func (s *GitAppsService) GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetBitbucketApp(ctx, id)
}

func (s *GitAppsService) SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveBitbucketApp(ctx, app)
}

func (s *GitAppsService) DeleteBitbucketApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteBitbucketApp(ctx, id)
}
