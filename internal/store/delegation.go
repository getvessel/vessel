package store

import (
	"context"

	"vessel.dev/vessel/internal/settings"
	"vessel.dev/vessel/internal/user"
)

// GetServerSettings delegates to the modular settings repository.
func (s *Store) GetServerSettings() (*settings.ServerSettings, error) {
	return settings.NewSQLiteRepository(s.db).GetServerSettings(context.Background())
}

// UpdateServerSettings delegates to the modular settings repository.
func (s *Store) UpdateServerSettings(cfg *settings.ServerSettings) error {
	return settings.NewSQLiteRepository(s.db).UpdateServerSettings(context.Background(), cfg)
}

// CreateUser delegates to the modular user repository.
func (s *Store) CreateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).CreateUser(context.Background(), u)
}

// GetUserByEmail delegates to the modular user repository.
func (s *Store) GetUserByEmail(email string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByEmail(context.Background(), email)
}

// GetUserByID delegates to the modular user repository.
func (s *Store) GetUserByID(id string) (*user.User, error) {
	return user.NewSQLiteRepository(s.db).GetUserByID(context.Background(), id)
}

// ListUsers delegates to the modular user repository.
func (s *Store) ListUsers() ([]user.User, error) {
	return user.NewSQLiteRepository(s.db).ListUsers(context.Background())
}

// UpdateUser delegates to the modular user repository.
func (s *Store) UpdateUser(u *user.User) error {
	return user.NewSQLiteRepository(s.db).UpdateUser(context.Background(), u)
}

// CreatePersonalAccessToken delegates to the modular user repository.
func (s *Store) CreatePersonalAccessToken(pat *user.PersonalAccessToken) error {
	return user.NewSQLiteRepository(s.db).CreatePAT(context.Background(), pat)
}

// ListPersonalAccessTokens delegates to the modular user repository.
func (s *Store) ListPersonalAccessTokens(userID string) ([]*user.PersonalAccessToken, error) {
	return user.NewSQLiteRepository(s.db).ListPATs(context.Background(), userID)
}

// DeletePersonalAccessToken delegates to the modular user repository.
func (s *Store) DeletePersonalAccessToken(id, userID string) error {
	return user.NewSQLiteRepository(s.db).DeletePAT(context.Background(), id, userID)
}
