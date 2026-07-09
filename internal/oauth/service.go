package oauth

import (
	"context"
	"errors"
	"strings"

	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/user"
)

// Service handles all OAuth provider management and 2FA logic.
type Service struct {
	repo         Repository
	userRepo     user.Repository
	tokenService *services.TokenService
}

func NewService(repo Repository, userRepo user.Repository, tokenService *services.TokenService) *Service {
	return &Service{repo: repo, userRepo: userRepo, tokenService: tokenService}
}

func (s *Service) ListProviders(ctx context.Context) ([]Provider, error) {
	return s.repo.ListProviders(ctx)
}

func (s *Service) SaveProvider(ctx context.Context, p *Provider) error {
	if p.ID == "" && p.ProviderName != "" {
		p.ID = strings.ToLower(p.ProviderName)
	}
	return s.repo.SaveProvider(ctx, p)
}

// HandleCallback exchanges the OAuth code for a user email, upserts the user, and returns a JWT token.
func (s *Service) HandleCallback(ctx context.Context, providerName, code string) (string, *user.User, error) {
	p, err := s.repo.GetProvider(ctx, providerName)
	if err != nil || p == nil {
		return "", nil, errors.New("oauth provider not found: " + providerName)
	}

	email, err := ExchangeCode(p, code)
	if err != nil || email == "" {
		return "", nil, errors.New("failed oauth code exchange")
	}

	u, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if u == nil {
		u = &user.User{
			Email:         email,
			PasswordHash:  "oauth-login-no-password",
			Role:          "member",
			OAuthProvider: p.ProviderName,
		}
		if err := s.userRepo.CreateUser(ctx, u); err != nil {
			return "", nil, errors.New("failed to create user account from oauth: " + err.Error())
		}
	}

	token, err := s.tokenService.GenerateToken(u)
	if err != nil {
		return "", nil, errors.New("failed generating token")
	}

	u.PasswordHash = ""
	return token, u, nil
}

func (s *Service) GetProvider(ctx context.Context, idOrName string) (*Provider, error) {
	return s.repo.GetProvider(ctx, idOrName)
}

// Setup2FA generates a new TOTP secret and recovery codes for a user.
func (s *Service) Setup2FA(ctx context.Context, userID, email string) (*TwoFASetupResponse, error) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		return nil, errors.New("failed generating totp secret")
	}
	recoveryCodes, err := GenerateRecoveryCodes(8)
	if err != nil {
		return nil, errors.New("failed generating recovery codes")
	}
	if err := s.repo.UpdateUserTOTP(ctx, userID, false, secret, recoveryCodes); err != nil {
		return nil, err
	}
	return &TwoFASetupResponse{
		Secret:        secret,
		QRCodeURI:     GenerateTOTPQRUri(email, secret),
		RecoveryCodes: recoveryCodes,
	}, nil
}

// Verify2FA validates the TOTP passcode and enables 2FA on success.
func (s *Service) Verify2FA(ctx context.Context, userID, passcode string) error {
	secret, recoveryCodes, err := s.repo.GetUserTOTPSecret(ctx, userID)
	if err != nil || secret == "" {
		return errors.New("totp setup has not been initiated for this user")
	}
	if !ValidateTOTP(secret, passcode) {
		return errors.New("invalid 6-digit totp verification code")
	}
	return s.repo.UpdateUserTOTP(ctx, userID, true, secret, recoveryCodes)
}

// Disable2FA removes TOTP for a user.
func (s *Service) Disable2FA(ctx context.Context, userID string) error {
	return s.repo.UpdateUserTOTP(ctx, userID, false, "", nil)
}
