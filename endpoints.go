package x_clone_post_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetEndpoint  endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
	PostEndpoint endpoint.Endpoint
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

type postRequest struct {
	Post Post
}

type postResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPostRequest)
		p, e := s.Get(ctx, req.ID)
		return getPostResponse{Post: p, Err: e}, nil
	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.GetList(ctx)
		return getListResponse{Posts: p, Err: e}, nil
	}
}

func MakePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		p, e := s.Post(ctx, req.Post)
		return postResponse{Post: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetEndpoint:  MakeGetEndpoint(s),
		ListEndpoint: MakeListEndpoint(s),
		PostEndpoint: MakePostEndpoint(s),
	}
}
