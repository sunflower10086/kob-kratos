package bot

import (
	"context"

	v1 "kob-kratos/api/gen/backend/v1"
	"kob-kratos/app/backend/internal/biz"

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
	err := s.botUc.AddBot(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetBotList(ctx context.Context, req *v1.GetBotListRequest) (*v1.GetBotListResponse, error) {
	return s.botUc.GetBotList(ctx, req)
}

func (s *Service) UpdateBot(ctx context.Context, req *v1.UpdateBotRequest) (*emptypb.Empty, error) {
	err := s.botUc.UpdateBot(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteBot(ctx context.Context, req *v1.DeleteBotRequest) (*emptypb.Empty, error) {
	err := s.botUc.DeleteBot(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
