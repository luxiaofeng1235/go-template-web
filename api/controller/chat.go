/*
 * @Description: 聊天控制器，处理聊天相关接口
 * @Version: 1.0.0
 * @Author: red
 * @Date: 2025-01-10
 * @LastEditTime: 2025-01-10
 */

package controller

import (
	"go-web-template/internal/models"
	"go-web-template/internal/service/api"
	"go-web-template/utils"
	"os"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/joho/godotenv"
)

type ChatController struct{}

func (c *ChatController) GetTokens(r *ghttp.Request) {
	err := godotenv.Load()
	if err != nil {
		utils.FailEncrypt(r, err, "读取失败")
		return
	}
	// 直接从.env环境变量获取AI聊天Token
	token := os.Getenv("AICHAT_SK_TOKEN")
	if token == "" {
		utils.FailEncrypt(r, nil, "获取失败")
		return
	}
	utils.Success(r, map[string]interface{}{
		"token": token,
	}, "获取Token成功")
}

// GetDeviceList 获取设备列表（对应PHP的chat/getDeviceList）
func (c *ChatController) GetDeviceList(r *ghttp.Request) {
	var req models.ChatUserListReq
	if err := r.Parse(&req); err != nil {
		utils.Fail(r, err, "参数解析失败")
		return
	}

	list, _, err := api.GetChatUserList(&req)
	if err != nil {
		utils.Fail(r, err, "获取设备列表失败")
		return
	}

	utils.Success(r, list, "用户列表")
}

// GetChatHistory 获取聊天历史记录（对应PHP的Chat/getChatHistory）
func (c *ChatController) GetChatHistory(r *ghttp.Request) {
	var req models.ChatHistoryReq
	if err := r.Parse(&req); err != nil {
		utils.Fail(r, err, "参数解析失败")
		return
	}

	list, _, err := api.GetChatHistoryByParams(&req)
	if err != nil {
		utils.Fail(r, err, "获取聊天历史失败")
		return
	}

	utils.Success(r, map[string]interface{}{
		"messages": list,
	}, "聊天历史")
}
