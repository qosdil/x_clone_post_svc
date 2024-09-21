package x_clone_post_srv

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetPost(ctx context.Context, id string) (post Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetPost", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetPost(ctx, id)
}

func (mw loggingMiddleware) GetPosts(ctx context.Context) (posts []Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetPosts", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetPosts(ctx)
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
