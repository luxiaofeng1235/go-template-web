package router

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// InitAdminRoutes 初始化管理后台路由（类似go-novel的admin_routes.InitRoutes）
func InitAdminRoutes(s *ghttp.Server) {
	// 创建管理后台控制器实例
	// TODO: 添加管理后台控制器
	
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
		
		// 用户管理
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
		
		// 产品管理
		productGroup := adminGroup.Group("/product")
		{
			productGroup.GET("/list", func(r *ghttp.Request) {
				r.Response.WriteJsonExit(map[string]interface{}{
					"code": 0,
					"msg":  "success",
					"data": []map[string]interface{}{
						{"id": 1, "name": "产品1", "price": 99.99},
						{"id": 2, "name": "产品2", "price": 199.99},
					},
				})
			})
		}
	}
}