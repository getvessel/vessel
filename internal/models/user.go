package models

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Role          string    `json:"role"`
	TOTPEnabled   bool      `json:"totpEnabled"`
	OAuthProvider string    `json:"oauthProvider,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type UserClaims struct {
	UserID      string `json:"sub"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	TOTPEnabled bool   `json:"totpEnabled"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	Prefix    string    `json:"prefix"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
}

type CreatePATRequest struct {
	Name string `json:"name"`
}

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
}
