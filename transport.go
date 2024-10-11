package x_clone_post_svc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"x_clone_post_svc/configs"
	model "x_clone_post_svc/model"
	service "x_clone_post_svc/service"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

const (
	httpStatusUnauthorizedMessage = "Unauthorized"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

type errorer interface {
	error() error
}

func codeFrom(err error) int {
	switch err {
	case model.ErrNotFound:
		return http.StatusNotFound
	case model.ErrAlreadyExists, model.ErrBadIDFormat:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getByIDRequest{ID: id}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeCreateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req createRequest

	// Extract the validated JWT user ID from auth middleware
	userID, _ := ctx.Value("user_id").(string)

	req.UserID = userID
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

type errField struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	code := err.Error()
	message, ok := model.Errors[err.Error()]

	// Set Code and Message with HTTP default statuses if not found in the map
	if !ok {
		message = strings.ToLower(http.StatusText(http.StatusInternalServerError))
		code = strings.ReplaceAll(message, " ", "_")
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errField{Code: code, Message: message},
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHTTPHandler(s service.Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}
	pathPrefix := "/posts"
	v1Path := "/v1" + pathPrefix
	r.Methods("GET").Path(v1Path + "/{id}").Handler(httptransport.NewServer(
		e.GetByIDEndpoint,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path(v1Path).Handler(httptransport.NewServer(
		e.ListEndpoint,
		decodeListRequest,
		encodeResponse,
		options...,
	))
	r.Handle(v1Path, jwtAuthMiddleware(configs.GetEnv("JWT_SECRET"))(httptransport.NewServer(
		e.CreateEndpoint,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))).Methods("POST")
	return r
}
