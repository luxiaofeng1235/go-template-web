package router

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"go-web-template/api/controller"
)

// InitAPIRoutes 初始化API路由（类似go-novel的api_routes.InitApiRoutes）
func InitAPIRoutes(s *ghttp.Server) {
	// 创建控制器实例
	userCtrl := &controller.UserController{}
	productCtrl := &controller.ProductController{}

	// API路由组
	apiGroup := s.Group("/api")
	{
		// 用户相关路由
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", userCtrl.Register)
			userGroup.POST("/login", userCtrl.Login)
			userGroup.GET("/profile", userCtrl.GetProfile)
		}

		// 产品相关路由
		productGroup := apiGroup.Group("/product")
		{
			productGroup.GET("/list", productCtrl.GetList)
			productGroup.GET("/detail", productCtrl.GetDetail)
			productGroup.POST("/create", productCtrl.Create)
		}
	}

	// 初始化聊天路由
	InitChatRoutes(s)
}
