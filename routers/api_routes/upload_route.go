/*
 * @file: upload_route.go
 * @description: 文件上传路由
 * @author: Claude
 * @created: 2025-09-09
 * @version: 1.0.0
 * @license: MIT License
 */
package api_routes

import (
	"go-web-template/api/controller"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitUploadRoutes 初始化API文件上传路由
func InitUploadRoutes(apiGroup *ghttp.RouterGroup) {
	uploadCtrl := &controller.UploadController{}

	// 直接在file路径下定义接口
	apiGroup.POST("/file/formimage", uploadCtrl.FormImage) // 图片上传接口
	apiGroup.POST("/file/formvideo", uploadCtrl.FormVideo) // 视频上传接口
}
