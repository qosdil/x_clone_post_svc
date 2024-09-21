package x_clone_post_srv

import "context"

type Post struct {
	ID   string `json:"id"`
	Post string `json:"post"`
}

type Service interface {
	GetPosts(ctx context.Context) ([]Post, error)
}

type dbService struct{}

func NewDbService() Service {
	return &dbService{}
}

func (s *dbService) GetPosts(ctx context.Context) (posts []Post, err error) {
	posts = append(posts, Post{ID: "fffa333", Post: "Hello"})
	posts = append(posts, Post{ID: "fffa334", Post: "Hi! This is my first post."})
	return posts, nil
}
