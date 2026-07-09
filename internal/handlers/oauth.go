package handlers

import (
	"github.com/labstack/echo/v4"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(s *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: s}
}

func (h *OAuthHandler) ListProviders(c echo.Context) error {
	providers, err := h.oauthService.ListProviders(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, providers)
}

func (h *OAuthHandler) SaveProvider(c echo.Context) error {
	var p models.OAuthProviderConfig
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.oauthService.SaveProvider(r.Context(), &p); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, p)
}

func (h *OAuthHandler) OAuthRedirect(c echo.Context) error {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}
	p, err := h.oauthService.GetProvider(r.Context(), providerName)
	if err != nil || p == nil {
		WriteError(w, http.StatusNotFound, "oauth provider not found or not enabled: "+providerName)
		return nil
	}
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to generate secure state token")
		return nil
	}
	state := hex.EncodeToString(stateBytes)
	authURL, err := services.GetAuthorizationURL(p, state)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) OAuthCallback(c echo.Context) error {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")
	code := c.QueryParam("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "missing authorization code parameter")
		return nil
	}
	token, _, err := h.oauthService.HandleCallback(r.Context(), providerName, code)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return nil
	}
	SetAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Setup2FA(c echo.Context) error {
	claims := ExtractClaims(r)
	if claims == nil || claims.UserID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	res, err := h.oauthService.Setup2FA(r.Context(), claims.UserID, claims.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, res)
}

func (h *OAuthHandler) Verify2FA(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	var payload struct {
		Passcode string `json:"passcode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Passcode == "" {
		WriteError(w, http.StatusBadRequest, "missing 6-digit passcode")
		return nil
	}
	if err := h.oauthService.Verify2FA(r.Context(), userID, payload.Passcode); err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "totp_enabled"})
}

func (h *OAuthHandler) Disable2FA(c echo.Context) error {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return nil
	}
	if err := h.oauthService.Disable2FA(r.Context(), userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "totp_disabled"})
}
