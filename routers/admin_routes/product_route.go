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
	"github.com/gogf/gf/v2/net/ghttp"
)

// InitProductRoutes 初始化管理后台产品路由
func InitProductRoutes(adminGroup *ghttp.RouterGroup) {
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
