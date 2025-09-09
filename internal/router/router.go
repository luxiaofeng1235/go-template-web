package router

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// InitAPI 初始化API路由（调用独立的路由初始化函数）
func InitAPI(s *ghttp.Server) {
	// 初始化API路由
	InitAPIRoutes(s)
	
	// 静态资源
	s.SetServerRoot("public")
	s.AddStaticPath("/static", "public/static")
}

// InitAdmin 初始化管理后台路由（调用独立的路由初始化函数）
func InitAdmin(s *ghttp.Server) {
	// 初始化管理后台路由
	InitAdminRoutes(s)
	
	// 静态资源
	s.SetServerRoot("public")
	s.AddStaticPath("/static", "public/static")
}