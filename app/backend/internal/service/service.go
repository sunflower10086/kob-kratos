package service

import (
	"kob-kratos/internal/service/bot"
	"kob-kratos/internal/service/rank"
	"kob-kratos/internal/service/record"
	"kob-kratos/internal/service/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	bot.NewService,
	record.NewService,
	rank.NewService,
	user.NewService,
)
