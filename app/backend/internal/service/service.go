package service

import (
	"kob-kratos/app/backend/internal/service/bot"
	"kob-kratos/app/backend/internal/service/rank"
	"kob-kratos/app/backend/internal/service/record"
	"kob-kratos/app/backend/internal/service/user"
	"kob-kratos/app/backend/internal/service/ws"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	bot.NewService,
	rank.NewService,
	record.NewService,
	user.NewService,
	ws.NewHub,
)
