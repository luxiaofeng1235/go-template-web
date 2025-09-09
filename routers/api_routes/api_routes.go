/*
 * @file: api.go
 * @description: api路由定义
 * @author: red <513072539@qq.com>
 * @created: 2025-09-05
 * @version: 1.0.0
 * @license: MIT License
 * @copyright: Copyright (c) 2025 red. All rights reserved.
 */
package api_routes

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// InitRoutes 初始化API所有路由
func InitRoutes(s *ghttp.Server) {
	ctx := gctx.New()

	// 添加请求日志中间件
	s.Use(func(r *ghttp.Request) {
		// 请求开始日志
		g.Log().Infof(r.Context(), "[API] %s %s - 来源IP: %s", r.Method, r.URL.Path, r.GetClientIp())

		// 继续处理请求
		r.Middleware.Next()

		// 请求结束日志
		g.Log().Infof(r.Context(), "[API] %s %s - 响应状态: %d", r.Method, r.URL.Path, r.Response.Status)
	})

	// API路由组
	apiGroup := s.Group("/api")
	{
		// g.Log().Info(ctx, "[API] 注册用户路由...")
		// InitUserRoutes(apiGroup)

		// g.Log().Info(ctx, "[API] 注册产品路由...")
		// InitProductRoutes(apiGroup)

		g.Log().Info(ctx, "[API] 注册文件管理路由...")
		InitFileRoutes(apiGroup)
		InitChatRoutes(apiGroup)
	}
}
