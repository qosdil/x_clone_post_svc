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
	Post PostResponse `json:"post,omitempty"`
	Err  error        `json:"err,omitempty"`
}

type listResponse struct {
	Posts []PostResponse `json:"posts,omitempty"`
	Err   error          `json:"err,omitempty"`
}

type postRequest struct {
	Post Post
}

type postResponse struct {
	Post PostResponse `json:"post,omitempty"`
	Err  error        `json:"err,omitempty"`
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
		p, e := s.GetList(ctx)
		return listResponse{Posts: p, Err: e}, nil
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
		GetByIDEndpoint: MakeGetByIDEndpoint(s),
		ListEndpoint:    MakeListEndpoint(s),
		PostEndpoint:    MakePostEndpoint(s),
	}
}
