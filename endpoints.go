package x_clone_post_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetPostEndpoint  endpoint.Endpoint
	GetListEndpoint  endpoint.Endpoint
	PostPostEndpoint endpoint.Endpoint
}

type getPostRequest struct {
	ID string
}

type getPostResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

type getListResponse struct {
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
		p, e := s.Get(ctx, req.ID)
		return getPostResponse{Post: p, Err: e}, nil
	}
}

func MakeGetListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetList(ctx)
		return getListResponse{Posts: p, Err: e}, nil
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
		GetListEndpoint:  MakeGetListEndpoint(s),
		PostPostEndpoint: MakePostPostEndpoint(s),
	}
}
