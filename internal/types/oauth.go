package types

import "time"

type OAuthProvider struct {
	ID           string    `json:"id"`
	ProviderName string    `json:"providerName"`
	Enabled      bool      `json:"enabled"`
	ClientID     string    `json:"clientId"`
	ClientSecret string    `json:"clientSecret"`
	RedirectURI  string    `json:"redirectUri"`
	BaseURL      string    `json:"baseUrl,omitempty"`
	Tenant       string    `json:"tenant,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type User2FA struct {
	UserID        string   `json:"userId"`
	TOTPSecret    string   `json:"-"`
	TOTPEnabled   bool     `json:"totpEnabled"`
	RecoveryCodes []string `json:"-"`
}

type TwoFASetupResponse struct {
	Secret        string   `json:"secret"`
	QRCodeURI     string   `json:"qrCodeUri"`
	RecoveryCodes []string `json:"recoveryCodes"`
}
