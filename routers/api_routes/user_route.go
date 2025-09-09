/*
 * @file: user_route.go
 * @description: 用户管理路由
 * @author: red <513072539@qq.com>
 * @created: 2025-09-05
 * @version: 1.0.0
 * @license: MIT License
 * @copyright: Copyright (c) 2025 red. All rights reserved.
 */
package api_routes

import (
	"go-web-template/api/controller"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitUserRoutes 初始化API用户路由
func InitUserRoutes(apiGroup *ghttp.RouterGroup) {
	userCtrl := &controller.UserController{}

	userGroup := apiGroup.Group("/user")
	{
		userGroup.POST("/register", userCtrl.Register)
		userGroup.POST("/login", userCtrl.Login)
		userGroup.GET("/profile", userCtrl.GetProfile)
	}
}
