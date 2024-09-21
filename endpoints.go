package x_clone_post_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetPostsEndpoint endpoint.Endpoint
}

type getPostsResponse struct {
	Posts []Post `json:"posts,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func MakeGetPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetPosts(ctx)
		return getPostsResponse{Posts: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetPostsEndpoint: MakeGetPostsEndpoint(s),
	}
}
