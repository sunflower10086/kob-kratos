package bot

import (
	"context"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *Service) AddBot(ctx context.Context, req *v1.AddBotRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) GetBotList(ctx context.Context, req *v1.GetBotListRequest) (*v1.GetBotListResponse, error) {
	return nil, nil
}

func (s *Service) UpdateBot(ctx context.Context, req *v1.UpdateBotRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Service) DeleteBot(ctx context.Context, req *v1.DeleteBotRequest) (*emptypb.Empty, error) {
	return nil, nil
}
