package x_clone_post_svc

import (
	"context"
	model "x_clone_post_svc/model"
	service "x_clone_post_svc/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateEndpoint  endpoint.Endpoint
	GetByIDEndpoint endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
}

type createRequest struct {
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}

type createResponse struct {
	Post model.Post `json:"post,omitempty"`
	Err  error      `json:"err,omitempty"`
}

type getByIDRequest struct {
	ID string
}

type getByIDResponse struct {
	Post model.Post `json:"post,omitempty"`
	Err  error      `json:"err,omitempty"`
}

type listResponse struct {
	Posts []model.Post `json:"posts,omitempty"`
	Err   error        `json:"err,omitempty"`
}

func MakeGetByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getByIDRequest)
		p, e := s.GetByID(ctx, req.ID)
		return getByIDResponse{Post: p, Err: e}, nil
	}
}

func MakeListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		p, e := s.List(ctx)
		return listResponse{Posts: p, Err: e}, nil
	}
}

func MakeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createRequest)
		p, e := s.Create(ctx, model.Post{
			Content: req.Content,
			User: model.User{
				ID: req.UserID,
			},
		})
		return createResponse{Post: p, Err: e}, nil
	}
}

func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateEndpoint:  MakeCreateEndpoint(s),
		GetByIDEndpoint: MakeGetByIDEndpoint(s),
		ListEndpoint:    MakeListEndpoint(s),
	}
}
