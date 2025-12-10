package service

import (
	"kob-kratos/app/backend/internal/service/bot"
	"kob-kratos/app/backend/internal/service/rank"
	"kob-kratos/app/backend/internal/service/record"
	"kob-kratos/app/backend/internal/service/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	bot.NewService,
	record.NewService,
	rank.NewService,
	user.NewService,
)
