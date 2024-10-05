package x_clone_post_svc

import (
	"context"
)

type Service interface {
	GetByID(ctx context.Context, id string) (Post, error)
	List(ctx context.Context) ([]Post, error)
	Post(ctx context.Context, post Post) (Post, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id string) (Post, error) {
	return s.repo.FirstByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Post, error) {
	return s.repo.Find(ctx)
}

// TODO Change to Create(), follows the conventions of X API
func (s *service) Post(ctx context.Context, post Post) (Post, error) {
	return s.repo.Create(ctx, post)
}
