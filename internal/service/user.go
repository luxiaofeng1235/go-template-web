package service

import (
	"context"
	"errors"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"

	"github.com/gogf/gf/v2/frame/g"
)

type userService struct{}

var User = userService{}

// Register 用户注册
func (s userService) Register(ctx context.Context, username, password, email string) error {
	// 检查用户名是否已存在
	count, err := g.Model(constant.TABLE_USER).Where("username", username).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.MSG_USER_EXIST)
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// 插入用户数据
	_, err = g.Model(constant.TABLE_USER).Data(g.Map{
		"username": username,
		"password": hashedPassword,
		"email":    email,
		"status":   constant.USER_STATUS_NORMAL,
	}).Insert()

	return err
}

// Login 用户登录
func (s userService) Login(ctx context.Context, username, password string) (string, error) {
	// 查询用户信息
	var user models.UserModel
	err := g.Model(constant.TABLE_USER).Where("username", username).Scan(&user)
	if err != nil {
		return "", err
	}
	if user.ID == 0 {
		return "", errors.New(constant.MSG_USER_NOT_FOUND)
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New(constant.MSG_PASSWORD_ERROR)
	}

	// 检查用户状态
	if user.Status != constant.USER_STATUS_NORMAL {
		return "", errors.New("账号已被禁用")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(g.Map{
		"user_id":  user.ID,
		"username": user.Username,
	})
	if err != nil {
		return "", err
	}

	// 缓存token到Redis
	_, err = g.Redis().Set(ctx, constant.REDIS_TOKEN_PREFIX+token, user.ID)
	if err != nil {
		g.Log().Warning(ctx, "Redis缓存token失败:", err)
	}

	return token, nil
}

// GetProfile 获取用户资料
func (s userService) GetProfile(ctx context.Context, userID int) (*models.UserProfileRes, error) {
	var user models.UserModel
	err := g.Model(constant.TABLE_USER).Where("id", userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New(constant.MSG_USER_NOT_FOUND)
	}

	profile := &models.UserProfileRes{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	return profile, nil
}

// GetUserByID 根据ID获取用户信息
func (s userService) GetUserByID(ctx context.Context, userID int) (*models.UserModel, error) {
	var user models.UserModel
	err := g.Model(constant.TABLE_USER).Where("id", userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New(constant.MSG_USER_NOT_FOUND)
	}
	return &user, nil
}
