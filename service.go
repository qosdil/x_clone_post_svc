package x_clone_post_srv

import (
	"context"
)

type Service interface {
	Get(ctx context.Context, id string) (Post, error)
	GetList(ctx context.Context) ([]Post, error)
	Post(ctx context.Context, post Post) (Post, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (post Post, err error) {
	post, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (s *service) GetList(ctx context.Context) (posts []Post, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Post(ctx context.Context, post Post) (Post, error) {
	post, err := s.repo.Create(ctx, post)
	if err != nil {
		return post, err
	}
	return post, nil
}
