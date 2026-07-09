package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/user"
)

type Service struct {
	userRepo     user.Repository
	settingsRepo settings.Repository
	tokenService *services.TokenService
}

func NewService(userRepo user.Repository, settingsRepo settings.Repository, tokenService *services.TokenService) *Service {
	return &Service{
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
		tokenService: tokenService,
	}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (*AuthResult, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	users, _ := s.userRepo.ListUsers(ctx)
	isInitial := len(users) == 0

	cfg, _ := s.settingsRepo.GetServerSettings(ctx)
	if cfg != nil && !cfg.RegistrationEnabled && !isInitial {
		return nil, errors.New("user registration is disabled on this server")
	}

	if cfg != nil && !isInitial && strings.TrimSpace(cfg.RegistrationDomainAllowlist) != "" {
		allowed := false
		for _, d := range strings.Split(cfg.RegistrationDomainAllowlist, ",") {
			d = strings.TrimSpace(d)
			if d != "" && strings.HasSuffix(strings.ToLower(req.Email), "@"+strings.ToLower(d)) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, errors.New("email domain is not allowed on this server")
		}
	}

	existing, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := "member"
	if isInitial {
		role = "admin"
	}

	u := &user.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	token, err := s.tokenService.GenerateToken(&models.User{ID: u.ID, Email: u.Email, Role: u.Role})
	if err != nil {
		return nil, errors.New("failed to issue authentication token")
	}

	u.PasswordHash = ""
	return &AuthResult{Token: token, User: u}, nil
}

func (s *Service) Signin(ctx context.Context, req SigninRequest) (*AuthResult, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	u, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || u == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.tokenService.GenerateToken(&models.User{ID: u.ID, Email: u.Email, Role: u.Role})
	if err != nil {
		return nil, errors.New("failed to issue authentication token")
	}

	u.PasswordHash = ""
	return &AuthResult{Token: token, User: u}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*user.User, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return nil, errors.New("user account not found")
	}
	u.PasswordHash = ""
	return u, nil
}
