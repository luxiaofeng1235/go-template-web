/*
 * @file: stream.go
 * @description: AI流式搜索服务 - 实现AI流式对话和搜索功能
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package service

import (
	"context"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/models"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

type StreamService struct{}

// ChatLog 聊天记录结构
type ChatLog struct {
	ID       int    `json:"id" gorm:"column:id"`
	ChatID   int    `json:"chat_id" gorm:"column:chat_id"`
	UserID   int    `json:"user_id" gorm:"column:user_id"`
	Model    int    `json:"model" gorm:"column:model"`
	Messages string `json:"messages" gorm:"column:messages"`
	Answer   string `json:"answer" gorm:"column:answer"`
	Status   int    `json:"status" gorm:"column:status"`
	Created  int64  `json:"created" gorm:"column:created"`
	Updated  int64  `json:"updated" gorm:"column:updated"`
}

// TableName 设置表名
func (ChatLog) TableName() string {
	return "ls_ai_chat_log"
}

// SendStream2 实现流式AI搜索功能
func (s *StreamService) SendStream2(ctx context.Context, r *ghttp.Request, req *models.AIStreamReq) error {
	// 设置SSE响应头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 获取或创建聊天记录
	chatLog, err := s.getChatLog(ctx, req)
	if err != nil {
		return fmt.Errorf("获取聊天记录失败: %v", err)
	}

	// 初始化AI模型配置
	modelConfig, err := s.initializeModel(req.Model)
	if err != nil {
		return fmt.Errorf("初始化AI模型失败: %v", err)
	}

	// 准备消息内容
	messages, err := s.prepareMessages(chatLog, req.Msg)
	if err != nil {
		return fmt.Errorf("准备消息失败: %v", err)
	}

	// 模拟流式响应 - 这里需要根据实际AI服务调用进行实现
	err = s.streamResponse(r, messages, modelConfig, chatLog.ID)
	if err != nil {
		return fmt.Errorf("流式响应失败: %v", err)
	}

	return nil
}

// getChatLog 获取或创建聊天记录
func (s *StreamService) getChatLog(ctx context.Context, req *models.AIStreamReq) (*ChatLog, error) {
	var chatLog ChatLog

	// 如果提供了chat_id，尝试查询现有记录
	if req.ChatId > 0 {
		err := global.DB.WithContext(ctx).Where("chat_id = ?", req.ChatId).First(&chatLog).Error
		if err == nil {
			// 找到现有记录
			if req.Restart == 1 {
				// 重新开始对话，清空之前的消息
				chatLog.Messages = ""
				chatLog.Answer = ""
				chatLog.Status = 0
				chatLog.Updated = time.Now().Unix()
				global.DB.WithContext(ctx).Save(&chatLog)
			}
			return &chatLog, nil
		}
	}

	// 创建新的聊天记录
	now := time.Now().Unix()
	chatLog = ChatLog{
		ChatID:   req.ChatId,
		UserID:   0, // 默认用户ID，可以从请求中获取
		Model:    req.Model,
		Messages: "",
		Answer:   "",
		Status:   0, // 0:进行中, 1:完成, 2:失败
		Created:  now,
		Updated:  now,
	}

	err := global.DB.WithContext(ctx).Create(&chatLog).Error
	if err != nil {
		return nil, fmt.Errorf("创建聊天记录失败: %v", err)
	}

	return &chatLog, nil
}

// initializeModel 初始化AI模型配置
func (s *StreamService) initializeModel(model int) (map[string]interface{}, error) {
	// 根据模型类型返回不同配置
	modelConfigs := map[int]map[string]interface{}{
		0: {
			"name":        "gpt-3.5-turbo",
			"temperature": 0.7,
			"max_tokens":  2048,
		},
		1: {
			"name":        "gpt-4",
			"temperature": 0.7,
			"max_tokens":  2048,
		},
		2: {
			"name":        "claude-3",
			"temperature": 0.7,
			"max_tokens":  2048,
		},
		3: {
			"name":        "custom-model",
			"temperature": 0.7,
			"max_tokens":  2048,
		},
	}

	config, exists := modelConfigs[model]
	if !exists {
		return nil, fmt.Errorf("不支持的模型类型: %d", model)
	}

	return config, nil
}

// prepareMessages 准备发送给AI的消息
func (s *StreamService) prepareMessages(chatLog *ChatLog, newMessage string) ([]map[string]string, error) {
	var messages []map[string]string

	// 添加系统提示
	messages = append(messages, map[string]string{
		"role":    "system",
		"content": "你是一个有用的AI助手，请根据用户的问题提供准确、有帮助的回答。",
	})

	// 解析历史消息（如果有）
	if chatLog.Messages != "" {
		// 这里应该解析JSON格式的历史消息
		// 为了简化，这里只添加当前消息
	}

	// 添加当前用户消息
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": newMessage,
	})

	return messages, nil
}

// streamResponse 流式响应处理
func (s *StreamService) streamResponse(r *ghttp.Request, messages []map[string]string, modelConfig map[string]interface{}, logID int) error {
	// 模拟AI流式响应 - 实际实现需要调用真实的AI服务
	responseText := "这是一个模拟的AI响应。在实际实现中，这里应该调用真实的AI服务API，如OpenAI、Claude等，并处理流式响应。"

	// 分片发送响应
	words := []string{"这是", "一个", "模拟的", "AI", "响应", "。", "在实际", "实现中", "，", "这里", "应该", "调用", "真实的", "AI", "服务", "API", "，", "如", "OpenAI", "、", "Claude", "等", "，", "并", "处理", "流式", "响应", "。"}

	fullResponse := ""
	for i, word := range words {
		// 模拟延迟
		time.Sleep(100 * time.Millisecond)

		fullResponse += word

		// 发送数据块
		data := fmt.Sprintf("data: %s\n\n", word)
		_, err := r.Response.Write([]byte(data))
		if err != nil {
			return fmt.Errorf("写入响应失败: %v", err)
		}
		r.Response.Flush()

		// 检查客户端是否断开连接
		if r.Context().Err() != nil {
			break
		}
	}

	// 更新聊天记录
	err := s.updateChatLog(r.Context(), logID, fullResponse, 1)
	if err != nil {
		global.Errlog.Error(r.Context(), "更新聊天记录失败: %v", err)
	}

	// 发送完成标识
	endData := fmt.Sprintf("data: [LOG_ID]:%d\n\ndata: [DONE]\n\n", logID)
	_, err = r.Response.Write([]byte(endData))
	if err != nil {
		return fmt.Errorf("发送完成标识失败: %v", err)
	}
	r.Response.Flush()

	return nil
}

// updateChatLog 更新聊天记录
func (s *StreamService) updateChatLog(ctx context.Context, logID int, answer string, status int) error {
	updateData := map[string]interface{}{
		"answer":  answer,
		"status":  status,
		"updated": time.Now().Unix(),
	}

	err := global.DB.WithContext(ctx).Model(&ChatLog{}).Where("id = ?", logID).Updates(updateData).Error
	if err != nil {
		return fmt.Errorf("更新数据库失败: %v", err)
	}

	return nil
}

// GetChatHistory 获取聊天历史记录
func (s *StreamService) GetChatHistory(ctx context.Context, chatID int, limit int) ([]*ChatLog, error) {
	var logs []*ChatLog

	query := global.DB.WithContext(ctx).Where("chat_id = ?", chatID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("created DESC").Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("查询聊天历史失败: %v", err)
	}

	return logs, nil
}

// DeleteChatLog 删除聊天记录
func (s *StreamService) DeleteChatLog(ctx context.Context, logID int) error {
	err := global.DB.WithContext(ctx).Delete(&ChatLog{}, logID).Error
	if err != nil {
		return fmt.Errorf("删除聊天记录失败: %v", err)
	}

	return nil
}
