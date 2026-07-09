package storage

import "context"

type Repository interface {
	Create(ctx context.Context, s *Storage) error
	GetByID(ctx context.Context, id string) (*Storage, error)
	List(ctx context.Context) ([]*Storage, error)
	Delete(ctx context.Context, id string) error
}
