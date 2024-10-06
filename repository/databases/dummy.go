package databases

import (
	"context"
	"errors"
	"fmt"
	"time"
	model "x_clone_post_svc/model"
	repository "x_clone_post_svc/repository"
)

var dummy struct {
	posts []model.Post
}

func NewDummyRepository() repository.Repository {
	var post model.Post
	for i := 1; i <= 3; i++ {
		post = model.Post{
			ID:        fmt.Sprintf("dummyPostID_%d", i),
			Content:   fmt.Sprintf("dummy content %d", i),
			CreatedAt: uint32(time.Now().Unix()),
			User: model.User{
				ID:       fmt.Sprintf("dummyUserID_%d", i),
				Username: fmt.Sprintf("dummyUsername_%d", i),
			},
		}
		dummy.posts = append(dummy.posts, post)
	}
	return &dummyRepository{}
}

type dummyRepository struct{}

func (r *dummyRepository) Create(ctx context.Context, post model.Post) (model.Post, error) {
	createdAt := uint32(time.Now().Unix())
	post.ID = "a dummy post ID, data not saved to any persistent storage"
	post.CreatedAt = createdAt
	post.User.Username = fmt.Sprintf("dummyUsername_%s", post.User.ID)
	return post, nil
}

func (r *dummyRepository) Find(ctx context.Context) (posts []model.Post, err error) {
	return dummy.posts, nil
}

func (r *dummyRepository) FirstByID(ctx context.Context, id string) (post model.Post, err error) {
	for _, post = range dummy.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return model.Post{}, errors.New("post not found")
}
