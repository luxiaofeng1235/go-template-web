/*
 * @file: ai_route.go
 * @description: AI服务路由定义
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package api_routes

import (
	controller "go-web-template/internal/controller/api"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitAiRoutes 初始化AI相关路由
func InitAiRoutes(group *ghttp.RouterGroup) {
	aiController := &controller.AiController{}

	// AI路由组
	aiGroup := group.Group("/ai")
	{
		// 图片生成相关接口
		aiGroup.POST("/to_image", aiController.ToImage)  // 生成图片
		aiGroup.GET("/get_image", aiController.GetImage) // 获取图片生成结果

		// 视频生成相关接口
		aiGroup.POST("/to_video", aiController.ToVideo)  // 生成视频
		aiGroup.GET("/get_video", aiController.GetVideo) // 获取视频生成结果

		// AI作品历史记录
		aiGroup.GET("/get_ai_work_list", aiController.GetAiWorkList) // 获取AI作品列表

		// AI流式搜索
		aiGroup.POST("/send_stream_2", aiController.SendStream2) // 流式AI搜索
	}
}
