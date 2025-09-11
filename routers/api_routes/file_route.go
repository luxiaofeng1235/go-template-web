/*
 * @file: file_route.go
 * @description: 文件管理路由
 * @author: Claude
 * @created: 2025-09-09
 * @version: 1.0.0
 * @license: MIT License
 */
package api_routes

import (
	controller "go-web-template/internal/controller/api"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitFileRoutes 初始化API文件路由
func InitFileRoutes(apiGroup *ghttp.RouterGroup) {
	fileCtrl := &controller.FileController{}

	// 直接在file路径下定义接口
	apiGroup.POST("/file/formImage", fileCtrl.FormImage) // 图片上传接口
	apiGroup.POST("/file/formVideo", fileCtrl.FormVideo) // 视频上传接口
}
