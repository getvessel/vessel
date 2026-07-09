package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
)

type claimsKeyType struct{}

var claimsKey = claimsKeyType{}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

func SetAuthCookie(w http.ResponseWriter, token string) {
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

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "vessel_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

func ExtractClaims(r *http.Request) *models.UserClaims {
	if c, ok := r.Context().Value(claimsKey).(*models.UserClaims); ok {
		return c
	}
	return nil
}

func ExtractUserID(r *http.Request) string {
	if c := ExtractClaims(r); c != nil {
		return c.UserID
	}
	return ""
}

func WithUserClaims(ctx context.Context, claims *models.UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}
