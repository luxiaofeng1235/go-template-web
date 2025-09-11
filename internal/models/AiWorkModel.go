package models

import (
	"time"

	"gorm.io/datatypes"
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
	CreateTime *time.Time     `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime *time.Time     `gorm:"column:update_time" json:"update_time"` // 更新时间
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
	CreateTime *time.Time             `form:"create_time" json:"create_time"` // 创建时间
	UpdateTime *time.Time             `form:"update_time" json:"update_time"` // 更新时间
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
	TaskID string                 `json:"task_id"`
	Params map[string]interface{} `json:"params"`
	Type   int8                   `json:"type"`
	Status int8                   `json:"status"`
	Work   interface{}            `json:"work"`
}

// AiWorkListRes AI工作任务列表响应
type AiWorkListRes struct {
	Page      string      `json:"page"`
	PageCount int         `json:"page_count"`
	Total     int64       `json:"total"`
	List      []AiWorkRes `json:"list"`
}

// CreateAiWorkRes 创建AI工作任务响应
type CreateAiWorkRes struct {
	TaskID    string     `json:"task_id"`
	Success   bool       `json:"success"`
	Timestamp *time.Time `json:"timestamp"`
}

// ToImageReq 生成图片请求参数
type ToImageReq struct {
	UserID    int    `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	Model     int    `json:"model" form:"model" v:"required|in:1,2#模型必须选择|模型参数错误：1高速模型，2细节模型"`
	Prompt    string `json:"prompt" form:"prompt" v:"required|max-length:500#提示词必须输入|提示词长度不能超过500字符"`
	Size      string `json:"size" form:"size" v:"required#图片尺寸必须选择"`
	N         int    `json:"n" form:"n" v:"required|between:1,4#生成数量必须输入|生成数量必须在1-4之间"`
	Watermark int    `json:"watermark" form:"watermark"` // 是否添加水印：1添加，2不添加（默认），可选参数
	Image     string `json:"image" form:"image"` // 图生图时传的参数，只能传单张，可选
}

// ToVideoReq 生成视频请求参数
type ToVideoReq struct {
	UserID string `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	To     int    `json:"to" form:"to" v:"required|in:1,2#视频类型必须选择|视频类型参数错误：1图生视频，2文生视频"`
	Prompt string `json:"prompt" form:"prompt" v:"required|max-length:500#提示词必须输入|提示词长度不能超过500字符"`
	ImgURL string `json:"img_url" form:"img_url"`
}

// GetTaskResultReq 获取任务结果的统一请求参数（图片和视频通用）
type GetTaskResultReq struct {
	TaskId string `json:"task_id" form:"task_id" v:"required#任务ID必须输入"`
	UserID string `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
}

// AIWorkListReq AI作品列表请求参数
type AIWorkListReq struct {
	UserID int `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	Page   int `json:"page" form:"page" v:"required|min:1#页码必须输入|页码必须大于等于1"`
	Type   int `json:"type" form:"type" v:"required|in:1,2,3#类型必须选择|类型参数错误"`
}

// AIStreamReq ai搜索的流式输出
type AIStreamReq struct {
	Model   int    `json:"model" form:"model" v:"required|between:0,3#模型名称必须输入|模型参数必须在0-3之间"`
	ChatId  int    `json:"chat_id" form:"chat_id"`
	Msg     string `json:"msg" form:"msg" v:"required#搜索内容必须输入"`
	Restart int    `json:"restart" form:"restart"`
}

// AIImageResult AI图片生成结果
type AIImageResult struct {
	Iamge []string `json:"iamge"` // 注意这里是拼写错误的iamge，保持和PHP版本一致
	Image []string `json:"image"`
}

// AIVideoResult AI视频生成结果
type AIVideoResult struct {
	Video []string `json:"video"`
}
