package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vessel.dev/vessel/internal/types"
)

type OAuthService struct{}

func NewOAuthService() *OAuthService {
	return &OAuthService{}
}

func (s *OAuthService) GetAuthorizationURL(p *types.OAuthProvider, state string) (string, error) {
	if !p.Enabled || p.ClientID == "" {
		return "", fmt.Errorf("oauth provider %s is not enabled or configured", p.ProviderName)
	}

	switch strings.ToLower(p.ProviderName) {
	case "github":
		return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "gitlab":
		baseURL := p.BaseURL
		if baseURL == "" {
			baseURL = "https://gitlab.com"
		}
		return fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=read_user+openid+profile+email&state=%s",
			strings.TrimRight(baseURL, "/"), url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "google":
		return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "azuread":
		tenant := p.Tenant
		if tenant == "" {
			tenant = "common"
		}
		return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			url.PathEscape(tenant), url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "discord":
		return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "authentik", "zitadel", "clerk", "infomaniak":
		if p.BaseURL == "" {
			return "", fmt.Errorf("base url is required for %s oauth", p.ProviderName)
		}
		authEndpoint := strings.TrimRight(p.BaseURL, "/") + "/oauth/authorize"
		if strings.ToLower(p.ProviderName) == "authentik" {
			authEndpoint = strings.TrimRight(p.BaseURL, "/") + "/application/o/authorize/"
		} else if strings.ToLower(p.ProviderName) == "zitadel" {
			authEndpoint = strings.TrimRight(p.BaseURL, "/") + "/oauth/v2/authorize"
		}
		return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			authEndpoint, url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "bitbucket":
		return fmt.Sprintf("https://bitbucket.org/site/oauth2/authorize?client_id=%s&response_type=code&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(state)), nil
	default:
		return "", fmt.Errorf("unsupported oauth provider: %s", p.ProviderName)
	}
}

func (s *OAuthService) ExchangeCode(p *types.OAuthProvider, code string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	switch strings.ToLower(p.ProviderName) {
	case "github":
		reqBody := map[string]string{
			"client_id":     p.ClientID,
			"client_secret": p.ClientSecret,
			"code":          code,
			"redirect_uri":  p.RedirectURI,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var tokenResp struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil || tokenResp.AccessToken == "" {
			return "", fmt.Errorf("failed to get github access token: %w", err)
		}

		userReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		userResp, err := client.Do(userReq)
		if err != nil {
			return "", err
		}
		defer userResp.Body.Close()

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err := json.NewDecoder(userResp.Body).Decode(&emails); err == nil {
			for _, e := range emails {
				if e.Primary {
					return e.Email, nil
				}
			}
			if len(emails) > 0 {
				return emails[0].Email, nil
			}
		}
		return "", fmt.Errorf("could not retrieve email from github")

	case "google", "gitlab", "azuread", "discord", "authentik", "zitadel":
		tokenURL := "https://oauth2.googleapis.com/token"
		userURL := "https://openidconnect.googleapis.com/v1/userinfo"

		if strings.ToLower(p.ProviderName) == "gitlab" {
			baseURL := p.BaseURL
			if baseURL == "" {
				baseURL = "https://gitlab.com"
			}
			tokenURL = strings.TrimRight(baseURL, "/") + "/oauth/token"
			userURL = strings.TrimRight(baseURL, "/") + "/api/v4/user"
		} else if strings.ToLower(p.ProviderName) == "discord" {
			tokenURL = "https://discord.com/api/oauth2/token"
			userURL = "https://discord.com/api/users/@me"
		} else if strings.ToLower(p.ProviderName) == "azuread" {
			tenant := p.Tenant
			if tenant == "" {
				tenant = "common"
			}
			tokenURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant)
			userURL = "https://graph.microsoft.com/oidc/userinfo"
		} else if strings.ToLower(p.ProviderName) == "authentik" {
			tokenURL = strings.TrimRight(p.BaseURL, "/") + "/application/o/token/"
			userURL = strings.TrimRight(p.BaseURL, "/") + "/application/o/userinfo/"
		} else if strings.ToLower(p.ProviderName) == "zitadel" {
			tokenURL = strings.TrimRight(p.BaseURL, "/") + "/oauth/v2/token"
			userURL = strings.TrimRight(p.BaseURL, "/") + "/oidc/v1/userinfo"
		}

		values := url.Values{
			"client_id":     {p.ClientID},
			"client_secret": {p.ClientSecret},
			"code":          {code},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {p.RedirectURI},
		}

		resp, err := client.PostForm(tokenURL, values)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var tokenResp struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil || tokenResp.AccessToken == "" {
			return "", fmt.Errorf("failed to exchange oauth code for token")
		}

		userReq, _ := http.NewRequest("GET", userURL, nil)
		userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		userResp, err := client.Do(userReq)
		if err != nil {
			return "", err
		}
		defer userResp.Body.Close()

		var userInfo struct {
			Email string `json:"email"`
		}
		bodyBytes, _ := io.ReadAll(userResp.Body)
		if err := json.Unmarshal(bodyBytes, &userInfo); err == nil && userInfo.Email != "" {
			return userInfo.Email, nil
		}

		return "", fmt.Errorf("could not extract email from provider user info")

	default:
		return "", fmt.Errorf("oauth exchange not implemented for provider: %s", p.ProviderName)
	}
}
