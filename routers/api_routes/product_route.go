package api_routes

import (
	controller "go-web-template/internal/controller/api"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitProductRoutes 初始化商品相关路由
func InitProductRoutes(group *ghttp.RouterGroup) {
	productController := &controller.ProductController{}

	// 商品路由组
	productGroup := group.Group("/product")
	{
		// GET /api/product/getProductList - 获取商品列表
		productGroup.GET("/getProductList", productController.GetProductList)
		
		// GET /api/product/getCategoryList - 获取分类列表  
		productGroup.GET("/getCategoryList", productController.GetCategoryList)
	}
}