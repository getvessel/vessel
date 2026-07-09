package oauth

import "time"

type Provider struct {
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

type TwoFASetupResponse struct {
	Secret        string   `json:"secret"`
	QRCodeURI     string   `json:"qrCodeUri"`
	RecoveryCodes []string `json:"recoveryCodes"`
}
