package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, userID, newEmail string) (*User, error) {
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed fetching user profile: %w", err)
	}
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}

	if strings.TrimSpace(newEmail) != "" {
		u.Email = strings.TrimSpace(newEmail)
	}
	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return nil, fmt.Errorf("failed updating profile: %w", err)
	}
	return u, nil
}

func (s *Service) ListPATs(ctx context.Context, userID string) ([]*PersonalAccessToken, error) {
	return s.repo.ListPATs(ctx, userID)
}

func (s *Service) CreatePAT(ctx context.Context, userID, name string) (*PersonalAccessToken, string, error) {
	if strings.TrimSpace(name) == "" {
		name = "Personal Access Token"
	}

	rawBytes := make([]byte, 24)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate random token: %w", err)
	}
	rawToken := fmt.Sprintf("vsl_user_%s", hex.EncodeToString(rawBytes))
	tokenHash := sha256.Sum256([]byte(rawToken))
	hashStr := hex.EncodeToString(tokenHash[:])

	pat := &PersonalAccessToken{
		UserID:    userID,
		Name:      name,
		TokenHash: hashStr,
		Prefix:    "vsl_user_",
	}

	if err := s.repo.CreatePAT(ctx, pat); err != nil {
		return nil, "", fmt.Errorf("failed to create PAT: %w", err)
	}
	return pat, rawToken, nil
}

func (s *Service) DeletePAT(ctx context.Context, id, userID string) error {
	return s.repo.DeletePAT(ctx, id, userID)
}
