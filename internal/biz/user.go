package biz

import (
	"context"
	"time"

	v1 "kob-kratos/api/backend/v1"
	"kob-kratos/internal/conf"
	"kob-kratos/internal/data/gormgen/query"
	"kob-kratos/internal/pkg/jwtc"
	"kob-kratos/pkg/codex"
	"kob-kratos/pkg/errx"

	"github.com/go-kratos/kratos/v2/log"
)

// User 用户实体
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
	Rating   int64  `json:"rating"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Insert(ctx context.Context, tx *query.Query, user *User) error
	Update(ctx context.Context, tx *query.Query, user *User) error
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userID int64) (*User, error)
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	Transaction(ctx context.Context, fn func(tx *query.Query) error) error
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepository

	jwtConf *conf.Jwt

	log *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepository, logger log.Logger, jwtConf *conf.Jwt) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		log:     log.NewHelper(logger),
		jwtConf: jwtConf,
	}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// 验证密码确认
	if req.Password != req.ConfirmedPassword {
		err := errx.New(codex.CodeConfirmPasswordError, "密码和确认密码不匹配")
		return nil, err
	}

	// 检查用户是否已存在
	existingUser, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		err := errx.Internal(err, "检查用户是否存在失败")
		return nil, err
	}
	if existingUser != nil {
		err := errx.New(codex.CodeUserExist, "用户名已存在")
		return nil, err
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
		err := errx.Internal(err, "用户注册失败")
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
		err := errx.Internal(err, "用户登录失败")
		return nil, err
	}

	if user == nil {
		return nil, errx.New(codex.CodeUserNotExist, "用户不存在")
	}

	// 这里应该验证密码，简化处理
	if user.Password != req.Password {
		return &v1.LoginResponse{
			Token: "",
		}, nil
	}

	// 生成JWT token
	token, err := generateJWTToken(uc.jwtConf.AccessSecret, user.ID)
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

// GetUserInfo 获取用户信息
func (uc *UserUsecase) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	user, err := uc.repo.GetUserInfo(ctx, req.UserId)
	if err != nil {
		uc.log.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, errx.New(codex.CodeUserNotExist, "用户不存在")
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
func generateJWTToken(jwtSk string, userID int64) (string, error) {
	payload := jwtc.Payload{
		Uid: userID,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwtc.GenJwtToken(jwtSk, &payload)
	if err != nil {
		err := errx.Internal(err, "生成JWT token失败")
		return "", err
	}
	return token, nil
}
