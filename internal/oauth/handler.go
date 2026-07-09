package oauth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	service       *Service
	extractClaims func(r *http.Request) (userID, email string)
}

func NewHandler(service *Service, extractClaims func(r *http.Request) (userID, email string)) *Handler {
	return &Handler{service: service, extractClaims: extractClaims}
}

func (h *Handler) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.ListProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, providers)
}

func (h *Handler) SaveProvider(w http.ResponseWriter, r *http.Request) {
	var p Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.service.SaveProvider(r.Context(), &p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	if idx := strings.Index(providerName, "/"); idx != -1 {
		providerName = providerName[:idx]
	}

	p, err := h.service.GetProvider(r.Context(), providerName)
	if err != nil || p == nil {
		writeError(w, http.StatusNotFound, "oauth provider not found or not enabled: "+providerName)
		return
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate secure state token")
		return
	}
	state := hex.EncodeToString(stateBytes)

	authURL, err := GetAuthorizationURL(p, state)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	providerName := strings.TrimPrefix(r.URL.Path, "/api/auth/oauth/")
	providerName = strings.TrimSuffix(providerName, "/callback")

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "missing authorization code parameter")
		return
	}

	token, _, err := h.service.HandleCallback(r.Context(), providerName, code)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	setAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) Setup2FA(w http.ResponseWriter, r *http.Request) {
	userID, email := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	res, err := h.service.Setup2FA(r.Context(), userID, email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	var payload struct {
		Passcode string `json:"passcode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Passcode == "" {
		writeError(w, http.StatusBadRequest, "missing 6-digit passcode")
		return
	}
	if err := h.service.Verify2FA(r.Context(), userID, payload.Passcode); err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "totp_enabled"})
}

func (h *Handler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := h.extractClaims(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	if err := h.service.Disable2FA(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "totp_disabled"})
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vessel_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   72 * 3600,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
