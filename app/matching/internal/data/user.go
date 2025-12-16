package data

import (
	"context"

	"kob-kratos/app/matching/internal/data/gormgen/model"
	"kob-kratos/app/matching/internal/data/gormgen/query"

	"github.com/go-kratos/kratos/v2/log"
)

var _ UserRepository = (*userRepo)(nil)

type UserRepository interface {
	FindUserById(ctx context.Context, userId int64) (*model.User, error)
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(data *Data, logger log.Logger) UserRepository {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

// Transaction 执行事务操作
func (r *userRepo) Transaction(ctx context.Context, fn func(tx *query.Query) error) error {
	return r.data.DB.Transaction(fn)
}

func (r *userRepo) FindUserById(ctx context.Context, userId int64) (*model.User, error) {
	userModel := r.data.DB.User
	return userModel.WithContext(ctx).Where(userModel.ID.Eq(userId)).First()
}
