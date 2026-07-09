package database

import "context"

type Repository interface {
	Create(ctx context.Context, db *Database) error
	GetByID(ctx context.Context, id string) (*Database, error)
	List(ctx context.Context) ([]*Database, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, db *Database) error
}