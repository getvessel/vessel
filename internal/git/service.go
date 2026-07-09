package git

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type AppService struct {
	ID            string
	ProjectID     string
	EnvironmentID string
	Name          string
	RepositoryURL string
	Branch        string
	ContainerID   string
}

type ProjectService interface {
	ListAppServicesByProject(projectID string) ([]*AppService, error)
}

type Service struct {
	repo       Repository
	httpClient *http.Client
	projects   ProjectService
}

func NewService(repo Repository, httpClient *http.Client) *Service {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Service{repo: repo, httpClient: httpClient}
}

func (svc *Service) WithProjectService(ps ProjectService) {
	svc.projects = ps
}

func (svc *Service) SaveProvider(ctx context.Context, userID string, req *GitConnectRequest) (*GitProviderConfig, error) {
	switch req.Provider {
	case "github", "gitlab":
	default:
		return nil, errors.New("unsupported git provider; must be 'github' or 'gitlab'")
	}
	if req.AccessToken == "" {
		return nil, errors.New("access token is required")
	}

	gp := &GitProviderConfig{
		UserID:      userID,
		Provider:    req.Provider,
		AccessToken: req.AccessToken,
		AccountName: req.AccountName,
	}
	if err := svc.repo.SaveProvider(ctx, gp); err != nil {
		return nil, fmt.Errorf("failed to save git provider: %w", err)
	}
	gp.AccessToken = ""
	return gp, nil
}

func (svc *Service) GetConnectedProviders(ctx context.Context, userID string) ([]map[string]any, error) {
	providers, err := svc.repo.ListProvidersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	providerMap := make(map[string]*GitProviderConfig)
	for _, gp := range providers {
		providerMap[gp.Provider] = gp
	}

	var results []map[string]any
	for _, provider := range []string{"github", "gitlab"} {
		if gp, ok := providerMap[provider]; ok && gp != nil {
			results = append(results, map[string]any{
				"provider":    provider,
				"connected":   true,
				"accountName": gp.AccountName,
				"updatedAt":   gp.UpdatedAt,
			})
		} else {
			results = append(results, map[string]any{
				"provider":  provider,
				"connected": false,
			})
		}
	}
	return results, nil
}

func (svc *Service) DisconnectProvider(ctx context.Context, userID, provider string) error {
	return svc.repo.DeleteProvider(ctx, userID, provider)
}

func (svc *Service) ListRepositories(ctx context.Context, userID, provider string) ([]GitRepository, error) {
	gp, err := svc.repo.GetProvider(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to load git credentials: %w", err)
	}
	if gp == nil || gp.AccessToken == "" {
		return nil, fmt.Errorf("user is not authenticated with %s", provider)
	}

	switch provider {
	case "github":
		return svc.listGitHubRepos(ctx, gp.AccessToken)
	case "gitlab":
		return svc.listGitLabRepos(ctx, gp.AccessToken)
	default:
		return nil, errors.New("unsupported provider: " + provider)
	}
}

func (svc *Service) listGitHubRepos(ctx context.Context, token string) ([]GitRepository, error) {
	reqURL := "https://api.github.com/user/repos?per_page=100&sort=updated"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := svc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, string(body))
	}

	var ghRepos []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		CloneURL string `json:"clone_url"`
		HTMLURL  string `json:"html_url"`
		Default  string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghRepos); err != nil {
		return nil, fmt.Errorf("failed to decode github repositories: %w", err)
	}

	var results []GitRepository
	for _, r := range ghRepos {
		results = append(results, GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

func (svc *Service) listGitLabRepos(ctx context.Context, token string) ([]GitRepository, error) {
	reqURL := "https://gitlab.com/api/v4/projects?membership=true&per_page=100&order_by=updated_at"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := svc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gitlab api returned status %d: %s", resp.StatusCode, string(body))
	}

	var glRepos []struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		FullName   string `json:"path_with_namespace"`
		Visibility string `json:"visibility"`
		CloneURL   string `json:"http_url_to_repo"`
		HTMLURL    string `json:"web_url"`
		Default    string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&glRepos); err != nil {
		return nil, fmt.Errorf("failed to decode gitlab projects: %w", err)
	}

	var results []GitRepository
	for _, r := range glRepos {
		results = append(results, GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Visibility == "private",
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

func (svc *Service) CloneOrPullRepository(ctx context.Context, projectID, targetDir string, logWriter io.Writer) error {
	if svc.projects == nil {
		return errors.New("project service not configured")
	}
	apps, err := svc.projects.ListAppServicesByProject(projectID)
	if err != nil || len(apps) == 0 {
		return errors.New("no application services found for project to checkout")
	}
	return svc.CloneOrPullAppRepository(ctx, apps[0], targetDir, logWriter)
}

func (svc *Service) CloneOrPullAppRepository(ctx context.Context, app *AppService, targetDir string, logWriter io.Writer) error {
	repoURL := strings.TrimSpace(app.RepositoryURL)
	if repoURL == "" {
		return errors.New("repositoryUrl is not set for service")
	}

	branch := strings.TrimSpace(app.Branch)
	if branch == "" {
		branch = "main"
	}

	authURL := svc.injectAuthTokenIfAvailable(ctx, repoURL)

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📥 [GitService] Preparing to sync codebase from %s (branch: %s)...\n", repoURL, branch)
	}

	gitDir := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🔄 [GitService] Existing local directory detected; pulling latest changes...\n")
		}
		fetchCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "fetch", "origin", branch)
		if out, err := fetchCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git fetch failed: %v (%s)", err, string(out))
		}
		resetCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "reset", "--hard", "origin/"+branch)
		if out, err := resetCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git reset failed: %v (%s)", err, string(out))
		}
		if logWriter != nil {
			fmt.Fprintf(logWriter, "✅ [GitService] Successfully updated local repository to latest commit on %s.\n", branch)
		}
		return nil
	}

	_ = os.RemoveAll(targetDir)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create build parent dir: %w", err)
	}

	cloneArgs := []string{"clone", "--depth", "1", "-b", branch, authURL, targetDir}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [GitService] Running git clone --depth 1 -b %s...\n", branch)
	}

	cloneCmd := exec.CommandContext(ctx, "git", cloneArgs...)
	var stderr bytes.Buffer
	cloneCmd.Stderr = &stderr
	if err := cloneCmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "Remote branch") && branch == "main" {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "⚠️ [GitService] Branch 'main' not found; retrying clone with repository default branch...\n")
			}
			_ = os.RemoveAll(targetDir)
			cloneCmd = exec.CommandContext(ctx, "git", "clone", "--depth", "1", authURL, targetDir)
			if errFallback := cloneCmd.Run(); errFallback != nil {
				return fmt.Errorf("git clone failed: %v (%s)", errFallback, stderr.String())
			}
		} else {
			return fmt.Errorf("git clone failed: %v (%s)", err, stderr.String())
		}
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [GitService] Successfully cloned repository into %s.\n", targetDir)
	}
	return nil
}

func (svc *Service) injectAuthTokenIfAvailable(ctx context.Context, repoURL string) string {
	u, err := url.Parse(repoURL)
	if err != nil || u.Scheme != "https" {
		return repoURL
	}

	var provider string
	if strings.Contains(u.Host, "github.com") {
		provider = "github"
	} else if strings.Contains(u.Host, "gitlab.com") {
		provider = "gitlab"
	} else {
		return repoURL
	}

	gp, err := svc.repo.GetProvider(ctx, "", provider)
	if err != nil || gp == nil || gp.AccessToken == "" {
		return repoURL
	}

	switch provider {
	case "github":
		u.User = url.UserPassword("x-access-token", gp.AccessToken)
	case "gitlab":
		u.User = url.UserPassword("oauth2", gp.AccessToken)
	}
	return u.String()
}
