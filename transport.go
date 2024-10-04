package x_clone_post_svc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"x_clone_post_svc/config"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	httpStatusUnauthorizedMessage = "Unauthorized"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrBadRouting      = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrNotFound        = errors.New("not found")
)

type errorer interface {
	error() error
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getRequest{ID: id}, nil
}

func decodeGetListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodePostRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req postRequest

	// Extract the validated JWT user ID from auth middleware
	userIDStr, _ := ctx.Value("user_id").(string)
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	req.Post.UserID = userID
	if e := json.NewDecoder(r.Body).Decode(&req.Post); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
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

// TODO move this functionality to API gateway x Auth svc
// jwtAuthMiddleware is a middleware to validate the JWT token
func jwtAuthMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, httpStatusUnauthorizedMessage, http.StatusUnauthorized)
				return
			}

			// Split the token to get the Bearer part
			tokenParts := strings.Split(tokenString, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, httpStatusUnauthorizedMessage, http.StatusUnauthorized)
				return
			}
			tokenString = tokenParts[1]

			// Parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, nil
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, httpStatusUnauthorizedMessage, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http.Error(w, httpStatusUnauthorizedMessage, http.StatusUnauthorized)
				return
			}

			// Set user ID in context
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}
	pathPrefix := "/posts"
	v1Path := "/v1" + pathPrefix
	r.Methods("GET").Path(v1Path + "/{id}").Handler(httptransport.NewServer(
		e.GetEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path(v1Path).Handler(httptransport.NewServer(
		e.ListEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))
	r.Handle(v1Path, jwtAuthMiddleware(config.GetEnv("JWT_SECRET"))(httptransport.NewServer(
		e.PostEndpoint,
		decodePostRequest,
		encodeResponse,
		options...,
	))).Methods("POST")
	return r
}
