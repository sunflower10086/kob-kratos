package user

import (
	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// Service is a user service.
type Service struct {
	v1.UnimplementedUserServiceServer

	log    *log.Helper
	userUc *biz.UserUsecase
}

func NewService(userUc *biz.UserUsecase, logger log.Logger) *Service {
	return &Service{
		log:    log.NewHelper(log.With(logger, "module", "service/user")),
		userUc: userUc,
	}
}
