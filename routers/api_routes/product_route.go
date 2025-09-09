/*
 * @file: product_route.go
 * @description: 产品路由
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

// InitProductRoutes 初始化API产品路由
func InitProductRoutes(apiGroup *ghttp.RouterGroup) {
	productCtrl := &controller.ProductController{}

	productGroup := apiGroup.Group("/product")
	{
		productGroup.GET("/list", productCtrl.GetList)
		productGroup.GET("/detail", productCtrl.GetDetail)
		productGroup.POST("/create", productCtrl.Create)
	}
}
