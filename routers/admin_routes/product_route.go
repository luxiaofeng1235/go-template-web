/*
 * @file: product_route.go
 * @description: 商品管理路由定义 - Admin后台管理模块
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package admin_routes

import (
	controller "go-web-template/internal/controller/admin"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitProductRoutes 初始化管理后台商品路由
func InitProductRoutes(adminGroup *ghttp.RouterGroup) {
	productController := &controller.ProductController{}
	fileController := &controller.FileController{}

	// 商品管理路由组
	productGroup := adminGroup.Group("/product")
	{
		productGroup.GET("/getProductList", productController.GetProductList)
		productGroup.GET("/getCategoryList", productController.GetCategoryList)
		productGroup.POST("/saveProduct", productController.SaveProduct)
		productGroup.POST("/deleteProduct", productController.DeleteProduct)
		// 图片上传路由
		productGroup.POST("/formImage", fileController.FormImage)
	}
}
