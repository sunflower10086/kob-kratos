package record

import (
	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// Service is a record service.
type Service struct {
	v1.UnimplementedRecordServiceServer

	recordUc *biz.RecordUsecase
	log      *log.Helper
}

func NewService(recordUc *biz.RecordUsecase, logger log.Logger) *Service {
	return &Service{
		recordUc: recordUc,
		log:      log.NewHelper(log.With(logger, "module", "service/record")),
	}
}
