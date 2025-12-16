package biz

import (
	matchpool "kob-kratos/app/matching/internal/biz/matchpool"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewMatchUsecase,
	matchpool.NewMatchPool,
)
