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

func (mw loggingMiddleware) Get(ctx context.Context, id string) (post Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Get", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Get(ctx, id)
}

func (mw loggingMiddleware) GetList(ctx context.Context) (posts []Post, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetList", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetList(ctx)
}

func (mw loggingMiddleware) Post(ctx context.Context, post Post) (Post, error) {
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
