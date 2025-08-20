package rank

import (
	"context"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// Service is a rank service.
type Service struct {
	v1.UnimplementedRankServiceServer

	rankUc *biz.RankUsecase
	log    *log.Helper
}

func NewService(rankUc *biz.RankUsecase, logger log.Logger) *Service {
	return &Service{
		rankUc: rankUc,
		log:    log.NewHelper(log.With(logger, "module", "service/rank")),
	}
}

func (s *Service) GetRankList(ctx context.Context, req *v1.GetRankListRequest) (*v1.GetRankListResponse, error) {
	return nil, nil
}
