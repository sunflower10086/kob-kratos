package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"kob-kratos/internal/biz"
	"kob-kratos/internal/data/gormgen/query"
)

var _ biz.RankRepository = (*rankRepo)(nil)

type rankRepo struct {
	data *Data
	log  *log.Helper
}

// NewRankRepository 创建排行榜仓储实例
func NewRankRepository(data *Data, logger log.Logger) biz.RankRepository {
	return &rankRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/rank")),
	}
}

// GetRankList 获取排行榜列表
func (r *rankRepo) GetRankList(ctx context.Context, page, pageSize int32) ([]*biz.RankUser, int64, error) {
	// 查询用户总数
	total, err := r.data.DB.WithContext(ctx).User.Count()
	if err != nil {
		return nil, 0, err
	}

	// 按评分降序查询用户列表
	modelUsers, err := r.data.DB.WithContext(ctx).User.
		Order(r.data.DB.User.Rating.Desc(), r.data.DB.User.ID.Asc()).
		Scopes(Paginate(int(page), int(pageSize))).
		Find()
	if err != nil {
		return nil, 0, err
	}

	// 转换为业务实体
	users := make([]*biz.RankUser, 0, len(modelUsers))
	for i, modelUser := range modelUsers {
		user := &biz.RankUser{
			ID:       modelUser.ID,
			Username: modelUser.Username,
			Rating:   modelUser.Rating,
			Number:   int32((page-1)*pageSize + int32(i) + 1), // 排名从1开始
		}

		if modelUser.Photo != nil {
			user.Photo = *modelUser.Photo
		}

		users = append(users, user)
	}

	return users, total, nil
}

// GetUserRank 获取用户排名
func (r *rankRepo) GetUserRank(ctx context.Context, userID int32) (*biz.RankUser, error) {
	// 获取用户信息
	modelUser, err := r.data.DB.WithContext(ctx).User.Where(r.data.DB.User.ID.Eq(userID)).First()
	if err != nil {
		r.log.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	// 计算用户排名（比该用户评分高的用户数量 + 1）
	higherRatingCount, err := r.data.DB.WithContext(ctx).User.
		Where(r.data.DB.User.Rating.Gt(modelUser.Rating)).
		Count()
	if err != nil {
		r.log.Errorf("计算用户排名失败: %v", err)
		return nil, err
	}

	// 如果评分相同，按ID升序排列
	sameRatingLowerIDCount, err := r.data.DB.WithContext(ctx).User.
		Where(r.data.DB.User.Rating.Eq(modelUser.Rating), r.data.DB.User.ID.Lt(modelUser.ID)).
		Count()
	if err != nil {
		r.log.Errorf("计算同评分用户排名失败: %v", err)
		return nil, err
	}

	rank := higherRatingCount + sameRatingLowerIDCount + 1

	user := &biz.RankUser{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Rating:   modelUser.Rating,
		Number:   int32(rank),
	}

	if modelUser.Photo != nil {
		user.Photo = *modelUser.Photo
	}

	return user, nil
}

// UpdateUserRating 更新用户评分（事务方法）
func (r *rankRepo) UpdateUserRating(ctx context.Context, tx *query.Query, userID int32, rating int32) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	result, err := db.WithContext(ctx).User.
		Where(db.User.ID.Eq(userID)).
		UpdateSimple(db.User.Rating.Value(rating))
	if err != nil {
		r.log.Errorf("更新用户评分失败: %v", err)
		return err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("用户不存在: userID=%d", userID)
		return err
	}

	r.log.Infof("成功更新用户评分: userID=%d, rating=%d", userID, rating)
	return nil
}

// Transaction 执行事务操作
func (r *rankRepo) Transaction(ctx context.Context, fn func(tx *query.Query) error) error {
	return r.data.DB.Transaction(fn)
}
