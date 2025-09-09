/*
 * @file: global.go
 * @description: 后台路由定义
 * @author: red <513072539@qq.com>
 * @created: 2025-09-01
 * @version: 1.0.0
 * @license: MIT License
 * @copyright: Copyright (c) 2025 red. All rights reserved.
 */

package admin_routes

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// InitRoutes 初始化管理后台所有路由
func InitRoutes(s *ghttp.Server) {
	ctx := gctx.New()

	// 添加请求日志中间件
	s.Use(func(r *ghttp.Request) {
		// 请求开始日志
		g.Log().Infof(r.Context(), "[ADMIN] %s %s - 来源IP: %s", r.Method, r.URL.Path, r.GetClientIp())

		// 继续处理请求
		r.Middleware.Next()

		// 请求结束日志
		g.Log().Infof(r.Context(), "[ADMIN] %s %s - 响应状态: %d", r.Method, r.URL.Path, r.Response.Status)
	})

	// 管理后台路由组
	adminGroup := s.Group("/admin")
	{
		// 管理后台首页
		adminGroup.GET("/", func(r *ghttp.Request) {
			r.Response.WriteJsonExit(map[string]interface{}{
				"code": 0,
				"msg":  "success",
				"data": "管理后台首页",
			})
		})

		g.Log().Info(ctx, "[ADMIN] 注册用户路由...")
		InitUserRoutes(adminGroup)

		g.Log().Info(ctx, "[ADMIN] 注册产品路由...")
		InitProductRoutes(adminGroup)
	}
}
