package x_clone_post_srv

import (
	"context"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	ID        string `json:"id"`
	Post      string `json:"post"`
	CreatedAt int64  `json:"created_at"`
	User      User   `json:"user"`
}

type Service interface {
	GetPosts(ctx context.Context) ([]Post, error)
}

type dbService struct{}

func NewDbService() Service {
	return &dbService{}
}

func (s *dbService) GetPosts(ctx context.Context) (posts []Post, err error) {
	posts = append(posts, Post{
		ID: "fffa333", Post: "Hello", CreatedAt: time.Now().Unix(),
		User: User{
			ID:       "fjjlp341k",
			Username: "mary",
		},
	})
	posts = append(posts, Post{ID: "fffa334", Post: "Hi! This is my first post.", CreatedAt: time.Now().Unix(),
		User: User{
			ID:       "j99fjjjf",
			Username: "alan388",
		}})
	return posts, nil
}
