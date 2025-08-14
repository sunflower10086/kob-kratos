package bot

import (
	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// Service is a bot service.
type Service struct {
	v1.UnimplementedBotServiceServer

	botUc *biz.BotUsecase
	log   *log.Helper
}

func NewService(botUc *biz.BotUsecase, logger log.Logger) *Service {
	return &Service{
		botUc: botUc,
		log:   log.NewHelper(log.With(logger, "module", "service/bot")),
	}
}
