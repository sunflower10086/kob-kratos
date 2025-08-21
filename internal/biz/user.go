package biz

import (
	"context"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/data/gormgen/query"

	"github.com/go-kratos/kratos/v2/log"
)

// User 用户实体
type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
	Rating   int32  `json:"rating"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Insert(ctx context.Context, tx *query.Query, user *User) error
	Update(ctx context.Context, tx *query.Query, user *User) error
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userID int32) (*User, error)
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	Transaction(ctx context.Context, fn func(tx *query.Query) error) error
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepository
	log  *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepository, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// 验证密码确认
	if req.Password != req.ConfirmedPassword {
		uc.log.Error("密码和确认密码不匹配")
		return &v1.RegisterResponse{
			Message: "密码和确认密码不匹配",
		}, nil
	}

	// 检查用户是否已存在
	existingUser, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		uc.log.Errorf("检查用户是否存在失败: %v", err)
		return nil, err
	}
	if existingUser != nil {
		return &v1.RegisterResponse{
			Message: "用户名已存在",
		}, nil
	}

	// 创建用户
	user := &User{
		Username: req.Username,
		Password: req.Password, // 密码会在data层进行加密
		Photo:    "",           // 默认头像
		Rating:   1500,         // 默认评分
	}

	err = uc.repo.Insert(ctx, nil, user)
	if err != nil {
		uc.log.Errorf("用户注册失败: %v", err)
		return nil, err
	}

	return &v1.RegisterResponse{
		Message: "注册成功",
	}, nil
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		uc.log.Errorf("用户登录失败: %v", err)
		return nil, err
	}

	if user == nil {
		return &v1.LoginResponse{
			Token: "",
		}, nil
	}

	// 这里应该验证密码，简化处理
	if user.Password != req.Password {
		return &v1.LoginResponse{
			Token: "",
		}, nil
	}

	// 生成JWT token
	token := generateJWTToken(user.ID, user.Username)

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

// GetUserInfo 获取用户信息
func (uc *UserUsecase) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	userID := parseStringToInt32(req.UserId)
	user, err := uc.repo.GetUserInfo(ctx, userID)
	if err != nil {
		uc.log.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	if user == nil {
		return &v1.GetUserInfoResponse{
			UserId:   0,
			Username: "",
			Photo:    "",
		}, nil
	}

	return &v1.GetUserInfoResponse{
		UserId:   user.ID,
		Username: user.Username,
		Photo:    user.Photo,
	}, nil
}

// GetUserByUsername 根据用户名获取用户
func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		uc.log.Errorf("根据用户名获取用户失败: %v", err)
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) error {
	err := uc.repo.Update(ctx, nil, user)
	if err != nil {
		uc.log.Errorf("更新用户信息失败: %v", err)
		return err
	}
	return nil
}

// generateJWTToken 生成JWT token（简化实现）
func generateJWTToken(userID int32, username string) string {
	return "fixed_jwt_token_for_development"
}
