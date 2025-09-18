package validate

import (
	"context"

	"kob-kratos/pkg/errx"

	"github.com/go-kratos/kratos/v2/middleware"
)

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err = v.Validate(); err != nil {
					return nil, errx.BadRequest(err, "请求参数错误")
				}
			}
			return handler(ctx, req)
		}
	}
}
