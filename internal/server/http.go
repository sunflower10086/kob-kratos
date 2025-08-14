package server

import (
	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/conf"
	"kob-kratos/internal/service/bot"
	"kob-kratos/internal/service/rank"
	"kob-kratos/internal/service/record"
	"kob-kratos/internal/service/user"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	user *user.Service,
	record *record.Service,
	rank *rank.Service,
	bot *bot.Service,
	logger log.Logger,
) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
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
	srv := http.NewServer(opts...)
	v1.RegisterUserServiceHTTPServer(srv, user)
	v1.RegisterRecordServiceHTTPServer(srv, record)
	v1.RegisterRankServiceHTTPServer(srv, rank)
	v1.RegisterBotServiceHTTPServer(srv, bot)
	return srv
}
