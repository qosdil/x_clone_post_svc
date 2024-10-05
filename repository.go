package x_clone_post_svc

import (
	"context"
)

// Repository follows gORM convention for the method namings
type Repository interface {
	Create(ctx context.Context, post Post) (Post, error)
	Find(ctx context.Context) ([]Post, error)
	FirstByID(ctx context.Context, id string) (Post, error)
}
