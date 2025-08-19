package biz

import (
	"context"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/data/gormgen/query"

	"github.com/go-kratos/kratos/v2/log"
)

// RankUser 排行榜用户实体
type RankUser struct {
	ID       int32  `json:"id"`
	Photo    string `json:"photo"`
	Username string `json:"username"`
	Rating   int32  `json:"rating"`
	Number   int32  `json:"number"`
}

// RankRepository 排行榜仓储接口
type RankRepository interface {
	// GetRankList 获取排行榜列表
	GetRankList(ctx context.Context, page, pageSize int32) ([]*RankUser, int64, error)
	// GetUserRank 获取用户排名
	GetUserRank(ctx context.Context, userID int32) (*RankUser, error)
	// UpdateUserRating 更新用户评分
	UpdateUserRating(ctx context.Context, tx *query.Query, userID int32, rating int32) error
	Transaction(ctx context.Context, fn func(tx *query.Query) error) error
}

// RankUsecase 排行榜用例
type RankUsecase struct {
	repo RankRepository
	log  *log.Helper
}

// NewRankUsecase 创建排行榜用例
func NewRankUsecase(logger log.Logger) *RankUsecase {
	return &RankUsecase{
		// repo: repo,
		log: log.NewHelper(logger),
	}
}

// GetRankList 获取排行榜列表
func (uc *RankUsecase) GetRankList(ctx context.Context, req *v1.GetRankListRequest) (*v1.GetRankListResponse, error) {
	// page := parseStringToInt32(req.Page)
	// if page <= 0 {
	// 	page = 1
	// }

	// users, userCount, err := uc.repo.GetRankList(ctx, page)
	// if err != nil {
	// 	uc.log.Errorf("获取排行榜列表失败: %v", err)
	// 	return nil, err
	// }

	// userList := make([]*v1.RankUser, 0, len(users))
	// for _, user := range users {
	// 	userList = append(userList, &v1.RankUser{
	// 		Id:       user.ID,
	// 		Photo:    user.Photo,
	// 		Username: user.Username,
	// 		Rating:   user.Rating,
	// 		Number:   user.Number,
	// 	})
	// }

	// return &v1.GetRankListResponse{
	// 	Users:     userList,
	// 	UserCount: userCount,
	// }, nil
	panic("implement me")
}

// GetUserRank 获取用户排名
func (uc *RankUsecase) GetUserRank(ctx context.Context, userID int32) (*RankUser, error) {
	// user, err := uc.repo.GetUserRank(ctx, userID)
	// if err != nil {
	// 	uc.log.Errorf("获取用户排名失败: %v", err)
	// 	return nil, err
	// }
	// return user, nil
	panic("implement me")
}

// UpdateUserRating 更新用户评分
func (uc *RankUsecase) UpdateUserRating(ctx context.Context, userID int32, rating int32) error {
	// err := uc.repo.UpdateUserRating(ctx, userID, rating)
	// if err != nil {
	// 	uc.log.Errorf("更新用户评分失败: %v", err)
	// 	return err
	// }
	// return nil
	panic("implement me")
}
