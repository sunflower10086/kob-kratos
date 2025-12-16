package biz

import (
	"context"
	"sync"

	pb "kob-kratos/api/gen/matching/v1"
	"kob-kratos/app/matching/internal/data"

	matchpool "kob-kratos/app/matching/internal/biz/matchpool"

	"github.com/go-kratos/kratos/v2/log"
)

type MatchUsecase struct {
	log    *log.Helper
	locker sync.Mutex

	userRepo  data.UserRepository
	matchPool *matchpool.MatchPool
}

// NewMatchUsecase 创建匹配用例
func NewMatchUsecase(logger log.Logger, userRepo data.UserRepository, pool *matchpool.MatchPool) *MatchUsecase {
	return &MatchUsecase{
		log:    log.NewHelper(log.With(logger, "module", "biz/match")),
		locker: sync.Mutex{},

		userRepo:  userRepo,
		matchPool: pool,
	}
}

func (uc *MatchUsecase) AddUser(ctx context.Context, req *pb.User) error {
	user, err := uc.userRepo.FindUserById(ctx, int64(req.UserId))
	if err != nil {
		return err
	}

	uc.locker.Lock()
	defer uc.locker.Unlock()

	player := &matchpool.Player{
		UserId: int32(user.ID),
		BotId:  int32(req.BotId),
		Rating: int32(user.Rating),
	}
	uc.matchPool.AddPlayer(player)
	uc.log.Infof("add user success %d %d", player.UserId, player.BotId)

	return nil
}

func (uc *MatchUsecase) RemoveUser(ctx context.Context, req *pb.User) error {
	uc.locker.Lock()
	defer uc.locker.Unlock()

	uc.matchPool.RemovePlayer(int32(req.UserId))
	uc.log.Infof("remove user success %d", req.UserId)

	return nil
}
