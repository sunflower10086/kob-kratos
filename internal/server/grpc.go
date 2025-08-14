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
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	user *user.Service,
	record *record.Service,
	rank *rank.Service,
	bot *bot.Service,
	logger log.Logger,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(srv, user)
	v1.RegisterRecordServiceServer(srv, record)
	v1.RegisterRankServiceServer(srv, rank)
	v1.RegisterBotServiceServer(srv, bot)
	return srv
}
