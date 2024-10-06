package repository

import (
	"context"
	"x_clone_post_svc/model"
)

// Repository follows gORM convention for the method namings
type Repository interface {
	Create(ctx context.Context, post model.Post) (model.Post, error)
	Find(ctx context.Context) ([]model.Post, error)
	FirstByID(ctx context.Context, id string) (model.Post, error)
}
