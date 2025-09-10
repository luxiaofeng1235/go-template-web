/*
 * @file: stream.go
 * @description: AI流式搜索服务 - 按照PHP sendStream2 逻辑实现
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"gorm.io/datatypes"
)

type StreamService struct{}

// AIMessage AI消息格式
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIModelConfig AI模型配置 - 对应PHP的模型配置
type AIModelConfig struct {
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Model       string            `json:"model"`
	System      string            `json:"system"`
	Temperature float64           `json:"temperature"`
	MaxTokens   int               `json:"max_tokens"`
	Stream      bool              `json:"stream"`
	Search      bool              `json:"search"` // 是否启用搜索功能
}

// SendStream2 实现流式AI搜索功能 - 按照PHP sendStream2逻辑
func (s *StreamService) SendStream2(r *ghttp.Request, req *models.AIStreamReq) error {
	// 设置SSE响应头
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 忽略用户中断
	ctx := r.Context()

	// 1. 获取/创建聊天记录 - 对应PHP getChatLog()
	chat, err := s.getChatLog(ctx, int64(req.ChatId), "0", int64(req.Model)) // 适配AiChatLog模型
	if err != nil {
		global.Errlog.Error(ctx, "获取聊天记录失败: %v", err)
		return fmt.Errorf("聊天记录不存在")
	}

	// 2. 初始化AI模型配置 - 对应PHP initializeModel()
	modelConfig, err := s.initializeModel(req.Model, 0) // is_deep_reflection暂时设为0
	if err != nil {
		global.Errlog.Error(ctx, "初始化AI模型失败: %v", err)
		return fmt.Errorf("初始化AI模型失败: %v", err)
	}

	// 3. 准备消息格式 - 对应PHP prepareMessages()
	messages, err := s.prepareMessages(chat, req.Msg, req.Restart, modelConfig.System)
	if err != nil {
		global.Errlog.Error(ctx, "准备消息失败: %v", err)
		return fmt.Errorf("准备消息失败: %v", err)
	}

	// 删除AI设置，保存用户消息 - 对应PHP的：$chatMessages = array_slice($messages, 1);
	if len(messages) > 1 {
		chatMessages := messages[1:] // 去掉系统提示
		chatJSON, _ := json.Marshal(chatMessages)
		chat.Chat = datatypes.JSON(chatJSON)
	}

	// 4. 准备JSON载荷和请求头
	jsonPayload, err := s.prepareJSONPayload(messages, modelConfig)
	if err != nil {
		return fmt.Errorf("准备JSON载荷失败: %v", err)
	}

	// 5. 调用第三方AI API - 对应PHP的异步HTTP请求
	err = s.callAIServiceWithStream(ctx, r, modelConfig, jsonPayload, chat)
	if err != nil {
		global.Errlog.Error(ctx, "调用AI服务失败: %v", err)
		return fmt.Errorf("调用AI服务失败: %v", err)
	}

	return nil
}

// getChatLog 获取或创建聊天记录 - 对应PHP getChatLog()
func (s *StreamService) getChatLog(ctx context.Context, chatID int64, userID string, modelID int64) (*models.AiChatLog, error) {
	var chat models.AiChatLog

	if chatID > 0 {
		// 查询现有聊天记录 - 按照PHP逻辑：where('id', $chat_id)->where('user_id', $uid)->where('model_id', $model)
		err := global.DB.WithContext(ctx).Where("id = ? AND user_id = ? AND model_id = ?", chatID, userID, modelID).First(&chat).Error
		if err == nil {
			return &chat, nil
		}
	}

	// 创建新的聊天记录 - 按照PHP逻辑创建新记录
	now := time.Now()
	chat = models.AiChatLog{
		UserID:     userID,
		ModelID:    modelID,
		Chat:       datatypes.JSON("[]"), // 对应PHP的 $chat->chat = '[]';
		CreateTime: &now,
		UpdateTime: &now,
	}

	err := global.DB.WithContext(ctx).Create(&chat).Error
	if err != nil {
		return nil, fmt.Errorf("创建聊天记录失败: %v", err)
	}

	return &chat, nil
}

// initializeModel 初始化AI模型配置 - 对应PHP initializeModel()
func (s *StreamService) initializeModel(model, isDeepReflection int) (*AIModelConfig, error) {
	// 使用constant中的配置
	modelConfig := constant.GetAIModelConfig(constant.AI_MODEL_SEARCH, isDeepReflection == 1)

	config := &AIModelConfig{
		URL:         constant.ALIYUN_CHAT_URL,
		Model:       modelConfig.Model,
		System:      modelConfig.System,
		Temperature: 0.3,
		MaxTokens:   32768,
		Stream:      true,
		Search:      modelConfig.Search,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + constant.ALIYUN_AI_API_KEY,
		},
	}

	return config, nil
}

// prepareMessages 准备发送给AI的消息 - 对应PHP prepareMessages()
func (s *StreamService) prepareMessages(chat *models.AiChatLog, newMessage string, restart int, systemPrompt string) ([]AIMessage, error) {
	var messages []AIMessage

	// 添加系统提示 - 使用动态传入的系统提示
	messages = append(messages, AIMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	// 如果不是重新开始，且有历史记录，则加载历史消息
	if restart != 1 && len(chat.Chat) > 0 && string(chat.Chat) != "[]" {
		var historyMessages []AIMessage
		err := json.Unmarshal(chat.Chat, &historyMessages)
		if err == nil {
			messages = append(messages, historyMessages...)
		}
	}

	// 添加当前用户消息
	messages = append(messages, AIMessage{
		Role:    "user",
		Content: newMessage,
	})

	return messages, nil
}

// prepareJSONPayload 准备发送给AI服务的JSON载荷
func (s *StreamService) prepareJSONPayload(messages []AIMessage, config *AIModelConfig) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"model":       config.Model,
		"messages":    messages,
		"temperature": config.Temperature,
		"max_tokens":  config.MaxTokens,
		"stream":      config.Stream,
	}

	return payload, nil
}

// callAIServiceWithStream 调用AI服务并处理流式响应 - 对应PHP的异步HTTP请求和流处理
func (s *StreamService) callAIServiceWithStream(ctx context.Context, r *ghttp.Request, config *AIModelConfig, payload map[string]interface{}, chat *models.AiChatLog) error {
	// 准备请求数据
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", config.URL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI服务返回错误状态码: %d", resp.StatusCode)
	}

	// 处理流式响应 - 对应PHP的逐行读取
	return s.handleStreamResponse(ctx, r, resp.Body, chat)
}

// handleStreamResponse 处理流式响应 - 对应PHP的Utils::readLine()和handleResponseData()
func (s *StreamService) handleStreamResponse(ctx context.Context, r *ghttp.Request, body io.Reader, chat *models.AiChatLog) error {
	scanner := bufio.NewScanner(body)
	fullContent := ""

	for scanner.Scan() {
		// 检查客户端是否断开连接 - 对应PHP的connection_aborted()
		select {
		case <-ctx.Done():
			return s.streamClose(r, chat, fullContent)
		default:
		}

		line := scanner.Text()

		// 处理SSE数据行
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// 检查是否为结束标记
			if data == "[DONE]" {
				break
			}

			// 解析AI响应并提取内容
			content := s.extractContentFromAIResponse(data)
			if content != "" {
				fullContent += content

				// 实时发送给客户端
				outputData := fmt.Sprintf("data: %s\n\n", content)
				r.Response.Write([]byte(outputData))
				r.Response.Flush()
			}
		}
	}

	// 流式响应结束
	return s.streamClose(r, chat, fullContent)
}


// extractContentFromAIResponse 从AI响应中提取文本内容
func (s *StreamService) extractContentFromAIResponse(data string) string {
	// 这里需要根据具体AI服务的响应格式来解析
	// 不同的AI服务（OpenAI、Claude等）响应格式不同
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return ""
	}

	// OpenAI格式示例
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if delta, ok := choice["delta"].(map[string]interface{}); ok {
				if content, ok := delta["content"].(string); ok {
					return content
				}
			}
		}
	}

	return ""
}

// streamClose 关闭流并保存聊天记录 - 对应PHP的streamClose()
func (s *StreamService) streamClose(r *ghttp.Request, chat *models.AiChatLog, content string) error {
	// 更新聊天记录，添加AI回复到聊天历史
	if content != "" {
		var historyMessages []AIMessage
		if len(chat.Chat) > 0 && string(chat.Chat) != "[]" {
			json.Unmarshal(chat.Chat, &historyMessages)
		}

		// 添加AI回复到历史
		historyMessages = append(historyMessages, AIMessage{
			Role:    "assistant",
			Content: content,
		})

		chatJSON, _ := json.Marshal(historyMessages)
		chat.Chat = datatypes.JSON(chatJSON)
	}

	now := time.Now()
	chat.UpdateTime = &now

	err := global.DB.Save(chat).Error
	if err != nil {
		global.Errlog.Error(context.Background(), "保存聊天记录失败: %v", err)
	}

	// 发送完成标识 - 匹配PHP格式: data: [LOG_ID]:ID 然后 data: [DONE]
	endData := fmt.Sprintf("data: [LOG_ID]:%d\n\ndata: [DONE]\n\n", chat.ID)
	r.Response.Write([]byte(endData))
	r.Response.Flush()

	return nil
}
