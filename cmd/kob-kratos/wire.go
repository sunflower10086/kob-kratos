//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kob-kratos/internal/biz"
	"kob-kratos/internal/conf"
	"kob-kratos/internal/data"
	"kob-kratos/internal/server"
	"kob-kratos/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, server.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
