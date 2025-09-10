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
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
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
		// 管理后台首页 - 默认指向产品管理页面
		adminGroup.GET("/", func(r *ghttp.Request) {
			r.Response.RedirectTo("/admin/product")
		})

		// 产品管理页面
		adminGroup.GET("/product", func(r *ghttp.Request) {
			// 读取layout模板
			layoutContent := gfile.GetContents("public/admin/components/layout.html")
			// 读取产品页面内容
			productContent := gfile.GetContents("public/admin/product/index.html")
			
			// 将产品内容插入到layout模板中
			finalContent := gstr.Replace(layoutContent, "{$content|raw}", productContent)
			
			r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			r.Response.Write(finalContent)
		})

		// 静态资源路由
		adminGroup.GET("/assets/*", func(r *ghttp.Request) {
			path := "public/admin" + r.URL.Path
			r.Response.ServeFile(path)
		})

		adminGroup.GET("/components/*", func(r *ghttp.Request) {
			path := "public/admin" + r.URL.Path  
			r.Response.ServeFile(path)
		})

	// admin-frontend静态资源路由
	s.Group("/admin-frontend").GET("/*", func(r *ghttp.Request) {
		path := r.URL.Path[1:] // 去掉开头的 /
		r.Response.ServeFile(path)
	})

		g.Log().Info(ctx, "[ADMIN] 注册用户路由...")
		InitUserRoutes(adminGroup)

		g.Log().Info(ctx, "[ADMIN] 注册产品路由...")
		InitProductRoutes(adminGroup)
	}
}
