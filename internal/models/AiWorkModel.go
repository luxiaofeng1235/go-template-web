package models

import (
	"gorm.io/datatypes"
	"time"
)

// AiWork AI工作任务主模型
type AiWork struct {
	ID         int64          `gorm:"column:id" json:"id"`                   // ID
	UserID     string         `gorm:"column:user_id" json:"user_id"`         // 用户ID
	TaskID     string         `gorm:"column:task_id" json:"task_id"`         // 任务ID
	Params     datatypes.JSON `gorm:"column:params" json:"params"`           // 参数
	Work       datatypes.JSON `gorm:"column:work" json:"work"`               // 工作内容
	Type       int8           `gorm:"column:type" json:"type"`               // 类型：1文本生成 2图片生成 3音频生成 4视频生成
	Status     int8           `gorm:"column:status" json:"status"`           // 状态：0待处理 1处理中 2已完成 3失败
	CreateTime time.Time      `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime time.Time      `gorm:"column:update_time" json:"update_time"` // 更新时间
}

func (*AiWork) TableName() string {
	return "1new9j_ai_to_work"
}

// 请求结构体

// CreateAiWorkReq 创建AI工作任务请求
type CreateAiWorkReq struct {
	UserID     string                 `form:"user_id" json:"user_id" `
	TaskID     string                 `form:"task_id" json:"task_id"`
	Params     map[string]interface{} `form:"params" json:"params" `
	Type       int8                   `form:"type" json:"type" `
	CreateTime time.Time              `form:"create_time" json:"create_time"` // 创建时间
	UpdateTime time.Time              `form:"update_time" json:"update_time"` // 更新时间
}

// UpdateAiWorkReq 更新AI工作任务请求
type UpdateAiWorkReq struct {
	TaskID string `form:"task_id" json:"task_id" `
	CreateAiWorkReq
	Work   map[string]interface{} `form:"work" json:"work"`
	Status int8                   `form:"status" json:"status" `
}

// GetAiWorkListReq 获取AI工作任务列表请求
type GetAiWorkListReq struct {
	UserID   string `form:"user_id" json:"user_id"`
	Type     int8   `form:"type" json:"type"`
	Status   int8   `form:"status" json:"status"`
	Page     int    `form:"page" json:"page" `
	PageSize int    `form:"page_size" json:"page_size" `
}

// GetAiWorkDetailReq 获取AI工作任务详情请求
type GetAiWorkDetailReq struct {
	TaskID string `form:"task_id" json:"task_id" `
	UserID string `form:"user_id" json:"user_id"`
}

// 响应结构体

// AiWorkRes AI工作任务响应
type AiWorkRes struct {
	ID         int64                  `json:"id"`
	UserID     string                 `json:"user_id"`
	TaskID     string                 `json:"task_id"`
	Params     map[string]interface{} `json:"params"`
	Work       map[string]interface{} `json:"work,omitempty"`
	Type       int8                   `json:"type"`
	TypeName   string                 `json:"type_name"`
	Status     int8                   `json:"status"`
	StatusName string                 `json:"status_name"`
	CreateTime time.Time              `json:"create_time"`
	UpdateTime time.Time              `json:"update_time"`
}

// AiWorkListRes AI工作任务列表响应
type AiWorkListRes struct {
	Works    []AiWorkRes `json:"works"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	HasMore  bool        `json:"has_more"`
}

// CreateAiWorkRes 创建AI工作任务响应
type CreateAiWorkRes struct {
	TaskID    string     `json:"task_id"`
	Success   bool       `json:"success"`
	Timestamp *time.Time `json:"timestamp"`
}

// 常量定义

// AI工作类型常量 - 与PHP版本保持一致
const (
	AiWorkTypeText  int8 = 1 // 文本生成
	AiWorkTypeImage int8 = 2 // 图片生成
	AiWorkTypeAudio int8 = 3 // 音频生成
	AiWorkTypeVideo int8 = 4 // 视频生成
)

// AI工作状态常量
const (
	AiWorkStatusPending    int8 = 0 // 待处理
	AiWorkStatusProcessing int8 = 1 // 处理中
	AiWorkStatusCompleted  int8 = 2 // 已完成
	AiWorkStatusFailed     int8 = 3 // 失败
)

// 获取类型名称
func GetAiWorkTypeName(workType int8) string {
	switch workType {
	case AiWorkTypeText:
		return "文本生成"
	case AiWorkTypeImage:
		return "图片生成"
	case AiWorkTypeAudio:
		return "音频生成"
	case AiWorkTypeVideo:
		return "视频生成"
	default:
		return "未知类型"
	}
}

// 获取状态名称
func GetAiWorkStatusName(status int8) string {
	switch status {
	case AiWorkStatusPending:
		return "待处理"
	case AiWorkStatusProcessing:
		return "处理中"
	case AiWorkStatusCompleted:
		return "已完成"
	case AiWorkStatusFailed:
		return "失败"
	default:
		return "未知状态"
	}
}
