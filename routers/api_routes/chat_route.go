/*
 * @Description: 聊天相关路由定义，处理聊天接口路由注册
 * @Version: 1.0.0
 * @Author: red
 * @Date: 2025-01-10
 * @LastEditTime: 2025-01-10
 */

package api_routes

import (
	"go-web-template/api/controller"

	"github.com/gogf/gf/v2/net/ghttp"
)

// InitChatRoutes 初始化聊天相关路由
func InitChatRoutes(apiGroup *ghttp.RouterGroup) {
	chatCtrl := &controller.ChatController{}

	// 聊天相关路由
	apiGroup.GET("/chat/getTokens", chatCtrl.GetTokens)         // 获取Token接口
	apiGroup.GET("/chat/getDeviceList", chatCtrl.GetDeviceList) // 获取设备列表接口

	// Chat相关路由（大写C开头，对应PHP的Chat控制器）
	apiGroup.GET("/chat/getChatHistory", chatCtrl.GetChatHistory) // 获取聊天历史记录接口
	apiGroup.POST("/chat/deviceAuth", chatCtrl.DeviceAuth)        //认证设备
	apiGroup.POST("/chat/saveUserData", chatCtrl.SaveUserData)    //保存用户资料
}
