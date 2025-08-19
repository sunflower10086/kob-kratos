package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kob-kratos/internal/biz"
	"kob-kratos/internal/data/gormgen/model"
	"kob-kratos/internal/data/gormgen/query"
)

var _ biz.UserRepository = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(data *Data, logger log.Logger) biz.UserRepository {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

// GetUserInfo 获取用户信息
func (r *userRepo) GetUserInfo(ctx context.Context, userID int32) (*biz.User, error) {
	modelUser, err := r.data.DB.WithContext(ctx).User.Where(r.data.DB.User.ID.Eq(userID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("用户不存在: userID=%d", userID)
			return nil, nil
		}
		r.log.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	return r.modelToBiz(modelUser), nil
}

// GetUserByUsername 根据用户名获取用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*biz.User, error) {
	modelUser, err := r.data.DB.WithContext(ctx).User.Where(r.data.DB.User.Username.Eq(username)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("用户不存在: username=%s", username)
			return nil, nil
		}
		r.log.Errorf("根据用户名获取用户失败: %v", err)
		return nil, err
	}

	return r.modelToBiz(modelUser), nil
}

// modelToBiz 将数据模型转换为业务实体
func (r *userRepo) modelToBiz(modelUser *model.User) *biz.User {
	user := &biz.User{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Password: modelUser.Password,
		Rating:   modelUser.Rating,
	}

	if modelUser.Photo != nil {
		user.Photo = *modelUser.Photo
	}

	return user
}

// Insert 插入用户（事务方法）
func (r *userRepo) Insert(ctx context.Context, tx *query.Query, user *biz.User) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		r.log.Errorf("密码加密失败: %v", err)
		return err
	}

	modelUser := &model.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Rating:   1500, // 默认评分
	}

	if user.Photo != "" {
		modelUser.Photo = &user.Photo
	}

	if err := db.WithContext(ctx).User.Create(modelUser); err != nil {
		r.log.Errorf("插入用户失败: %v", err)
		return err
	}

	// 更新业务实体的ID
	user.ID = modelUser.ID
	return nil
}

// Update 更新用户（事务方法）
func (r *userRepo) Update(ctx context.Context, tx *query.Query, user *biz.User) error {
	// 如果tx为空，使用r.data.DB
	db := tx
	if db == nil {
		db = r.data.DB
	}

	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if user.Photo != "" {
		updateData["photo"] = user.Photo
	}

	if user.Rating > 0 {
		updateData["rating"] = user.Rating
	}

	result, err := db.WithContext(ctx).User.
		Where(db.User.ID.Eq(user.ID)).
		Updates(updateData)
	if err != nil {
		r.log.Errorf("更新用户失败: %v", err)
		return err
	}

	if result.RowsAffected == 0 {
		r.log.Warnf("用户不存在: userID=%d", user.ID)
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Transaction 执行事务操作
func (r *userRepo) Transaction(ctx context.Context, fn func(tx *query.Query) error) error {
	return r.data.DB.Transaction(fn)
}

