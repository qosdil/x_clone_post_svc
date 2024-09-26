package x_clone_post_srv

import (
	"context"
)

type Service interface {
	GetPost(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
	PostPost(ctx context.Context, post Post) (err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPost(ctx context.Context, id string) (post Post, err error) {
	post, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (s *service) GetPosts(ctx context.Context) (posts []Post, err error) {
	return s.repo.Find(ctx)
}

func (s *service) PostPost(ctx context.Context, post Post) (err error) {
	return nil
}
