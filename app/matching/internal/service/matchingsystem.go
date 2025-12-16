package service

import (
	"context"

	pb "kob-kratos/api/gen/matching/v1"
	"kob-kratos/app/matching/internal/biz"

	"google.golang.org/protobuf/types/known/emptypb"
)

type MatchingSystemService struct {
	pb.UnimplementedMatchingSystemServer

	matchUc *biz.MatchUsecase
}

func NewMatchingSystemService(matchUc *biz.MatchUsecase) *MatchingSystemService {
	return &MatchingSystemService{
		matchUc: matchUc,
	}
}

func (s *MatchingSystemService) AddUser(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	err := s.matchUc.AddUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MatchingSystemService) Remove(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	err := s.matchUc.RemoveUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
