package models

import "go-web-template/internal/constant"

// UserSession 用户会话主模型
type UserSession struct {
	ID               int64  `gorm:"column:id" json:"id"`                                 // ID
	SessionID        string `gorm:"column:session_id" json:"session_id"`                 // 会话ID
	WebsocketFD      int64  `gorm:"column:websocket_fd" json:"websocket_fd"`             // WebSocket文件描述符
	AccessKey        string `gorm:"column:access_key" json:"access_key"`                 // 访问密钥
	FirstLoginTime   int64  `gorm:"column:first_login_time" json:"first_login_time"`     // 首次登录时间
	LastActivityTime int64  `gorm:"column:last_activity_time" json:"last_activity_time"` // 最后活动时间
	SessionData      string `gorm:"column:session_data" json:"session_data"`             // 会话数据
	ClientIP         string `gorm:"column:client_ip" json:"client_ip"`                   // 客户端IP
	UserAgent        string `gorm:"column:user_agent" json:"user_agent"`                 // 用户代理
	Status           int8   `gorm:"column:status" json:"status"`                         // 状态：1活跃 2离线 3超时
	CreatedAt        int64  `gorm:"column:created_at" json:"created_at"`                 // 创建时间
	UpdatedAt        int64  `gorm:"column:updated_at" json:"updated_at"`                 // 更新时间
}

func (*UserSession) TableName() string {
	return "ls_user_sessions"
}

// 请求结构体

// CreateSessionReq 创建会话请求
type CreateSessionReq struct {
	AccessKey   string `form:"access_key" json:"access_key" `
	ClientIP    string `form:"client_ip" json:"client_ip"`
	UserAgent   string `form:"user_agent" json:"user_agent"`
	WebsocketFD int64  `form:"websocket_fd" json:"websocket_fd"`
	SessionData string `form:"session_data" json:"session_data"`
	CreatedAt   int64  `form:"created_at" json:"created_at"` // 创建时间
	UpdatedAt   int64  `form:"updated_at" json:"updated_at"` // 更新时间
}

// UpdateSessionReq 更新会话请求
type UpdateSessionReq struct {
	SessionID string `form:"session_id" json:"session_id" `
	CreateSessionReq
}

// GetSessionListReq 获取会话列表请求
type GetSessionListReq struct {
	AccessKey string `form:"access_key" json:"access_key"`
	Status    int8   `form:"status" json:"status"`
	Page      int    `form:"page" json:"page" `
	PageSize  int    `form:"page_size" json:"page_size" `
	ClientIP  string `form:"client_ip" json:"client_ip"`
}

// GetSessionDetailReq 获取会话详情请求
type GetSessionDetailReq struct {
	SessionID string `form:"session_id" json:"session_id" `
	AccessKey string `form:"access_key" json:"access_key"`
}

// UpdateSessionStatusReq 更新会话状态请求
type UpdateSessionStatusReq struct {
	SessionID string `form:"session_id" json:"session_id" `
	Status    int8   `form:"status" json:"status" `
}

// 响应结构体

// UserSessionRes 用户会话响应
type UserSessionRes struct {
	ID               int64  `json:"id"`
	SessionID        string `json:"session_id"`
	WebsocketFD      int64  `json:"websocket_fd"`
	AccessKey        string `json:"access_key"`
	FirstLoginTime   int64  `json:"first_login_time"`
	LastActivityTime int64  `json:"last_activity_time"`
	SessionData      string `json:"session_data,omitempty"`
	ClientIP         string `json:"client_ip"`
	UserAgent        string `json:"user_agent"`
	Status           int8   `json:"status"`
	StatusName       string `json:"status_name"`
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
	Duration         int64  `json:"duration"`     // 会话持续时间（秒）
	IsActive         bool   `json:"is_active"`    // 是否活跃
	DeviceType       string `json:"device_type"`  // 设备类型（从UserAgent解析）
	BrowserType      string `json:"browser_type"` // 浏览器类型
}

// UserSessionListRes 用户会话列表响应
type UserSessionListRes struct {
	Sessions []UserSessionRes `json:"sessions"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	HasMore  bool             `json:"has_more"`
}

// CreateSessionRes 创建会话响应
type CreateSessionRes struct {
	SessionID string `json:"session_id"`
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
}

// OnlineStatsRes 在线统计响应
type OnlineStatsRes struct {
	TotalOnline    int64 `json:"total_online"`    // 总在线人数
	ActiveSessions int64 `json:"active_sessions"` // 活跃会话数
	WebsocketConns int64 `json:"websocket_conns"` // WebSocket连接数
	TodayLogins    int64 `json:"today_logins"`    // 今日登录数
}

// 获取状态名称
func GetSessionStatusName(status int8) string {
	switch status {
	case constant.SessionStatusActive:
		return "活跃"
	case constant.SessionStatusOffline:
		return "离线"
	case constant.SessionStatusTimeout:
		return "超时"
	default:
		return "未知状态"
	}
}
