package x_clone_post_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetPostEndpoint  endpoint.Endpoint
	GetPostsEndpoint endpoint.Endpoint
	PostPostEndpoint endpoint.Endpoint
}

type getPostRequest struct {
	ID string
}

type getPostResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

type getPostsResponse struct {
	Posts []Post `json:"posts,omitempty"`
	Err   error  `json:"err,omitempty"`
}

type postPostRequest struct {
	Post Post
}

type postPostResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

func MakeGetPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPostRequest)
		p, e := s.GetPost(ctx, req.ID)
		return getPostResponse{Post: p, Err: e}, nil
	}
}

func MakeGetPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetPosts(ctx)
		return getPostsResponse{Posts: p, Err: e}, nil
	}
}

func MakePostPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postPostRequest)
		p, e := s.PostPost(ctx, req.Post)
		return postPostResponse{Post: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetPostEndpoint:  MakeGetPostEndpoint(s),
		GetPostsEndpoint: MakeGetPostsEndpoint(s),
		PostPostEndpoint: MakePostPostEndpoint(s),
	}
}
