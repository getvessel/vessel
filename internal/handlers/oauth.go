package handlers

import (
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

func (h *OAuthHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.oauthService.ListProviders(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, providers)
}

func (h *OAuthHandler) SaveProvider(w http.ResponseWriter, r *http.Request) {
	var p models.OAuthProviderConfig
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.oauthService.SaveProvider(r.Context(), &p); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, p)
}

func (h *OAuthHandler) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}
	p, err := h.oauthService.GetProvider(r.Context(), providerName)
	if err != nil || p == nil {
		WriteError(w, http.StatusNotFound, "oauth provider not found or not enabled: "+providerName)
		return
	}
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to generate secure state token")
		return
	}
	state := hex.EncodeToString(stateBytes)
	authURL, err := services.GetAuthorizationURL(p, state)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")
	code := r.URL.Query().Get("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "missing authorization code parameter")
		return
	}
	token, _, err := h.oauthService.HandleCallback(r.Context(), providerName, code)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	SetAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	claims := ExtractClaims(r)
	if claims == nil || claims.UserID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	res, err := h.oauthService.Setup2FA(r.Context(), claims.UserID, claims.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, res)
}

func (h *OAuthHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	var payload struct {
		Passcode string `json:"passcode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Passcode == "" {
		WriteError(w, http.StatusBadRequest, "missing 6-digit passcode")
		return
	}
	if err := h.oauthService.Verify2FA(r.Context(), userID, payload.Passcode); err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "totp_enabled"})
}

func (h *OAuthHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	userID := ExtractUserID(r)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	if err := h.oauthService.Disable2FA(r.Context(), userID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "totp_disabled"})
}
