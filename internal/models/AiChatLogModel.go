package models

import (
	"time"

	"gorm.io/datatypes"
)

// AiChatLog AI聊天日志主模型
type AiChatLog struct {
	ID         int64          `gorm:"column:id" json:"id"`                   // ID
	UserID     string         `gorm:"column:user_id" json:"user_id"`         // 用户ID
	ModelID    int64          `gorm:"column:model_id" json:"model_id"`       // 模型ID
	Chat       datatypes.JSON `gorm:"column:chat" json:"chat"`               // 聊天内容
	CreateTime int64          `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime int64          `gorm:"column:update_time" json:"update_time"` // 更新时间
}

func (*AiChatLog) TableName() string {
	return "1new9j_ai_chat_log"
}

// 请求结构体

// CreateAiChatLogReq 创建AI聊天日志请求
type CreateAiChatLogReq struct {
	UserID     string                 `form:"user_id" json:"user_id" `
	ModelID    int64                  `form:"model_id" json:"model_id" `
	Chat       map[string]interface{} `form:"chat" json:"chat" `
	CreateTime int64                  `form:"create_time" json:"create_time"` // 创建时间
	UpdateTime int64                  `form:"update_time" json:"update_time"` // 更新时间
}

// GetAiChatLogListReq 获取AI聊天日志列表请求
type GetAiChatLogListReq struct {
	UserID   string `form:"user_id" json:"user_id"`
	ModelID  int64  `form:"model_id" json:"model_id"`
	Page     int    `form:"page" json:"page" `
	PageSize int    `form:"page_size" json:"page_size" `
}

// GetAiChatLogDetailReq 获取AI聊天日志详情请求
type GetAiChatLogDetailReq struct {
	ID     int64  `form:"id" json:"id" `
	UserID string `form:"user_id" json:"user_id"`
}

// UpdateAiChatLogReq 更新AI聊天日志请求
type UpdateAiChatLogReq struct {
	ID int64 `form:"id" json:"id" `
	CreateAiChatLogReq
}

// 响应结构体

// AiChatLogRes AI聊天日志响应
type AiChatLogRes struct {
	ID         int64                  `json:"id"`
	UserID     string                 `json:"user_id"`
	ModelID    int64                  `json:"model_id"`
	ModelName  string                 `json:"model_name,omitempty"`
	Chat       map[string]interface{} `json:"chat"`
	CreateTime int64                  `json:"create_time"`
	UpdateTime int64                  `json:"update_time"`
}

// AiChatLogListRes AI聊天日志列表响应
type AiChatLogListRes struct {
	ChatLogs []AiChatLogRes `json:"chat_logs"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	HasMore  bool           `json:"has_more"`
}

// CreateAiChatLogRes 创建AI聊天日志响应
type CreateAiChatLogRes struct {
	ID        int64      `json:"id"`
	Success   bool       `json:"success"`
	Timestamp *time.Time `json:"timestamp"`
}

// AiChatStats AI聊天统计响应
type AiChatStatsRes struct {
	TotalChats  int64            `json:"total_chats"`
	TodayChats  int64            `json:"today_chats"`
	ModelStats  map[string]int64 `json:"model_stats"`
	RecentChats []AiChatLogRes   `json:"recent_chats"`
}

// ChatSession 聊天会话结构
type ChatSession struct {
	SessionID   string        `json:"session_id"`         // 会话ID
	Title       string        `json:"title"`              // 会话标题
	Messages    []ChatMessage `json:"messages"`           // 消息列表
	TotalTokens int           `json:"total_tokens"`       // 总token使用量
	StartTime   time.Time     `json:"start_time"`         // 会话开始时间
	EndTime     *time.Time    `json:"end_time,omitempty"` // 会话结束时间
}
