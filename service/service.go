package service

import (
	"context"
	"errors"
	"x_clone_post_svc/model"
	"x_clone_post_svc/repository"
)

type Service interface {
	Create(ctx context.Context, post model.Post) (model.Post, error)
	GetByID(ctx context.Context, id string) (model.Post, error)
	List(ctx context.Context) ([]model.Post, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, post model.Post) (model.Post, error) {
	if post.Content == "" {
		return model.Post{}, errors.New("content cannot be empty")
	}
	return s.repo.Create(ctx, post)
}

func (s *service) GetByID(ctx context.Context, id string) (model.Post, error) {
	return s.repo.FirstByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]model.Post, error) {
	return s.repo.Find(ctx)
}
