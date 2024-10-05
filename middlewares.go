package x_clone_post_svc

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetByID(ctx context.Context, id string) (postResponse PostResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetByID", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetByID(ctx, id)
}

func (mw loggingMiddleware) List(ctx context.Context) (posts []PostResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "List", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.List(ctx)
}

func (mw loggingMiddleware) Post(ctx context.Context, post Post) (PostResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Post", "took", time.Since(begin), "err", nil)
	}(time.Now())
	return mw.next.Post(ctx, post)
}

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
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
