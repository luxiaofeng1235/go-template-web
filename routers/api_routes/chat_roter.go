/*
 * @file: file_route.go
 * @description: 聊天的主要路由
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

// InitFileRoutes 初始化API文件路由
func InitChatRoutes(apiGroup *ghttp.RouterGroup) {
	fileCtrl := &controller.ChatController{}
	// 直接在file路径下定义接口
	apiGroup.GET("/chat/getTokens", fileCtrl.GetTokens) // 图片上传接口
}
