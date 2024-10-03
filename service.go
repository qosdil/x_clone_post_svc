package x_clone_post_svc

import (
	"context"
)

type Service interface {
	Get(ctx context.Context, id string) (PostResponse, error)
	GetList(ctx context.Context) ([]PostResponse, error)
	Post(ctx context.Context, post Post) (PostResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (PostResponse, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID.Hex(),
		Content:   post.Content,
		CreatedAt: post.CreatedAt.T,
		User: User{
			ID: post.UserID.Hex(),
		},
	}, nil
}

func (s *service) GetList(ctx context.Context) (postReponses []PostResponse, err error) {
	return s.repo.Find(ctx)
}

func (s *service) Post(ctx context.Context, post Post) (PostResponse, error) {
	post, err := s.repo.Create(ctx, post)
	if err != nil {
		return PostResponse{}, err
	}
	return PostResponse{
		ID:        post.ID.Hex(),
		Content:   post.Content,
		CreatedAt: post.CreatedAt.T,
		User: User{
			ID: post.UserID.Hex(),
		},
	}, nil
}
