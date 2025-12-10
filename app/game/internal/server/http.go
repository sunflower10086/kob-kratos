package server

import (
	"kob-kratos/app/game/internal/conf"
	"kob-kratos/app/game/internal/pkg/middlewares"
	"kob-kratos/pkg/httpencoder"
	"kob-kratos/pkg/middlewares/validate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bc *conf.Bootstrap,
	logger log.Logger,
) *http.Server {
	c := bc.Server
	confJwt := bc.Jwt
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(recovery.Recovery(), tracing.Server()).Prefix("/api").Build(),
			middlewares.Jwt(confJwt.GetAccessSecret()),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	opts = append(opts, http.ResponseEncoder(httpencoder.SuccessEncoder))
	opts = append(opts, http.ErrorEncoder(httpencoder.ErrorEncoder))

	srv := http.NewServer(opts...)
	return srv
}
