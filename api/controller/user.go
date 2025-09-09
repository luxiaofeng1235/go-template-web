package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"go-web-template/internal/service"
	"go-web-template/utils"
)

type UserController struct{}

// Register 用户注册
func (c *UserController) Register(r *ghttp.Request) {
	username := r.Get("username").String()
	password := r.Get("password").String()
	email := r.Get("email").String()

	if username == "" || password == "" {
		utils.ParamError(r, "用户名和密码不能为空")
		return
	}

	err := service.User.Register(r.Context(), username, password, email)
	if err != nil {
		utils.Fail(r, err, "注册失败")
		return
	}

	utils.Success(r, nil, "注册成功")
}

// Login 用户登录
func (c *UserController) Login(r *ghttp.Request) {
	username := r.Get("username").String()
	password := r.Get("password").String()

	if username == "" || password == "" {
		utils.ParamError(r, "用户名和密码不能为空")
		return
	}

	token, err := service.User.Login(r.Context(), username, password)
	if err != nil {
		utils.FailEncrypt(r, err, "登录失败")
		return
	}

	// 使用加密响应返回token
	tokenData := utils.TokenData{
		Token: token,
	}
	utils.SuccessEncrypt(r, tokenData, "登录成功")
}

// GetProfile 获取用户信息
func (c *UserController) GetProfile(r *ghttp.Request) {
	userID := r.Get("user_id").Int()
	if userID <= 0 {
		utils.ParamError(r, "用户ID无效")
		return
	}

	profile, err := service.User.GetProfile(r.Context(), userID)
	if err != nil {
		utils.FailEncrypt(r, err, "获取用户信息失败")
		return
	}

	utils.SuccessEncrypt(r, profile, "获取成功")
}