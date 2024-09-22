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
	GetPost(ctx context.Context, id string) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
	PostPost(ctx context.Context, post Post) (err error)
}

type dbService struct{}

func NewDbService() Service {
	return &dbService{}
}

func (s *dbService) GetPost(ctx context.Context, id string) (post Post, err error) {
	post = Post{ID: id, Post: "Hi! This is my first post.", CreatedAt: time.Now().Unix(),
		User: User{
			ID:       "j99fjjjf",
			Username: "alan388",
		}}
	return post, nil
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

func (s *dbService) PostPost(ctx context.Context, post Post) (err error) {
	return nil
}
