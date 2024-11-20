package user

import "context"

// Storage is a special interface for all the databases, that will use on server.
// For example you can use mongodb or postgresql.
type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindAll(ctx context.Context) ([]User, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}