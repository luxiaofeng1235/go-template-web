/*
 * @file: user_rote.go
 * @description: 用户路由定义
 * @author: red <513072539@qq.com>
 * @created: 2025-09-01
 * @version: 1.0.0
 * @license: MIT License
 * @copyright: Copyright (c) 2025 red. All rights reserved.
 */
package admin_routes

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// InitUserRoutes 初始化管理后台用户路由
func InitUserRoutes(adminGroup *ghttp.RouterGroup) {
	userGroup := adminGroup.Group("/user")
	{
		userGroup.GET("/list", func(r *ghttp.Request) {
			r.Response.WriteJsonExit(map[string]interface{}{
				"code": 0,
				"msg":  "success",
				"data": []map[string]interface{}{
					{"id": 1, "username": "admin", "email": "admin@example.com"},
					{"id": 2, "username": "user1", "email": "user1@example.com"},
				},
			})
		})

		userGroup.GET("/info", func(r *ghttp.Request) {
			r.Response.WriteJsonExit(map[string]interface{}{
				"code": 0,
				"msg":  "success",
				"data": map[string]interface{}{
					"id":       1,
					"username": "admin",
					"email":    "admin@example.com",
					"status":   1,
				},
			})
		})
	}
}
