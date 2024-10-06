package service_test

import (
	"context"
	"errors"
	"testing"
	"x_clone_post_svc/model"
	"x_clone_post_svc/repository/databases"
	"x_clone_post_svc/service"

	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
	s   = service.NewService(databases.NewDummyRepository())
)

func TestCreate(t *testing.T) {
	// Empty content
	createdPost, err := s.Create(ctx, model.Post{})
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("content cannot be empty"), err)

	// Good request
	content := "Hello!"
	userID := "some_user_id"
	createdPost, err = s.Create(ctx, model.Post{
		Content: content,
		User: model.User{
			ID: userID,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, content, createdPost.Content)
	assert.Equal(t, userID, createdPost.User.ID)
}

func TestGetByID(t *testing.T) {
	// Non-existence post
	post, err := s.GetByID(ctx, "some_id")
	assert.NotNil(t, err)
	assert.Empty(t, post)

	// Good request
	post, err = s.GetByID(ctx, "dummyPostID_3")
	assert.Nil(t, err)
	assert.NotEmpty(t, post)
}

func TestList(t *testing.T) {
	posts, err := s.List(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, posts)
}
