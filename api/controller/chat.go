package controller

import (
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
