package user

import "context"

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, u *User) error

	CreatePAT(ctx context.Context, pat *PersonalAccessToken) error
	ListPATs(ctx context.Context, userID string) ([]*PersonalAccessToken, error)
	DeletePAT(ctx context.Context, id, userID string) error
}
