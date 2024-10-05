package x_clone_post_svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetByIDEndpoint endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
	PostEndpoint    endpoint.Endpoint
}

type getRequest struct {
	ID string
}

type getResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

type listResponse struct {
	Posts []Post `json:"posts,omitempty"`
	Err   error  `json:"err,omitempty"`
}

type postRequest struct {
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}

type postResponse struct {
	Post Post  `json:"post,omitempty"`
	Err  error `json:"err,omitempty"`
}

func MakeGetByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		p, e := s.GetByID(ctx, req.ID)
		return getResponse{Post: p, Err: e}, nil
	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.List(ctx)
		return listResponse{Posts: p, Err: e}, nil
	}
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		p, e := s.Create(ctx, Post{
			Content: req.Content,
			User: User{
				ID: req.UserID,
			},
		})
		return postResponse{Post: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetByIDEndpoint: MakeGetByIDEndpoint(s),
		ListEndpoint:    MakeListEndpoint(s),
		PostEndpoint:    MakeCreateEndpoint(s),
	}
}
