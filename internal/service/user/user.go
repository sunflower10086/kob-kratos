package user

import (
	"context"

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

func (s *Service) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	return s.userUc.Register(ctx, req)
}

func (s *Service) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	return nil, nil
}

func (s *Service) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	return nil, nil
}
