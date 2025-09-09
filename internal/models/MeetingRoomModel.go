package models

import (
	"gorm.io/datatypes"
)

// MeetingRoom 会议室主模型
type MeetingRoom struct {
	ID               int            `gorm:"column:id" json:"id"`                                 // ID
	MeetingID        string         `gorm:"column:meeting_id" json:"meeting_id"`                 // 会议ID
	CreatorAccessKey string         `gorm:"column:creator_access_key" json:"creator_access_key"` // 创建者访问密钥
	CallType         int8           `gorm:"column:call_type" json:"call_type"`                   // 通话类型：1语音 2视频
	Status           int8           `gorm:"column:status" json:"status"`                         // 状态：1进行中 2已结束 3已取消
	Participants     datatypes.JSON `gorm:"column:participants" json:"participants"`             // 参与者JSON数据
	CreatedAt        int64          `gorm:"column:created_at" json:"created_at"`                 // 创建时间
	UpdatedAt        int64          `gorm:"column:updated_at" json:"updated_at"`                 // 更新时间
}

func (*MeetingRoom) TableName() string {
	return "ls_meeting_rooms"
}

// 请求结构体

// CreateMeetingRoomReq 创建会议室请求
type CreateMeetingRoomReq struct {
	MeetingID        string                 `form:"meeting_id" json:"meeting_id"`
	CreatorAccessKey string                 `form:"creator_access_key" json:"creator_access_key"`
	CallType         int8                   `form:"call_type" json:"call_type"`
	Participants     map[string]interface{} `form:"participants" json:"participants"`
	CreatedAt        int64                  `form:"created_at" json:"created_at"` // 创建时间
	UpdatedAt        int64                  `form:"updated_at" json:"updated_at"` // 更新时间
}

// UpdateMeetingRoomReq 更新会议室请求
type UpdateMeetingRoomReq struct {
	CreateMeetingRoomReq
	ID     int  `form:"id" json:"id"`
	Status int8 `form:"status" json:"status"`
}

// GetMeetingRoomListReq 获取会议室列表请求
type GetMeetingRoomListReq struct {
	CreatorAccessKey string `form:"creator_access_key" json:"creator_access_key"`
	CallType         int8   `form:"call_type" json:"call_type"`
	Status           int8   `form:"status" json:"status"`
	Page             int    `form:"page" json:"page"`
	PageSize         int    `form:"page_size" json:"page_size"`
}

// GetMeetingRoomDetailReq 获取会议室详情请求
type GetMeetingRoomDetailReq struct {
	ID        int    `form:"id" json:"id"`
	MeetingID string `form:"meeting_id" json:"meeting_id"`
}

// JoinMeetingRoomReq 加入会议室请求
type JoinMeetingRoomReq struct {
	MeetingID string `form:"meeting_id" json:"meeting_id"`
	AccessKey string `form:"access_key" json:"access_key"`
	Nickname  string `form:"nickname" json:"nickname"`
}

// LeaveMeetingRoomReq 离开会议室请求
type LeaveMeetingRoomReq struct {
	MeetingID string `form:"meeting_id" json:"meeting_id"`
	AccessKey string `form:"access_key" json:"access_key"`
}

// 响应结构体

// MeetingRoomRes 会议室响应
type MeetingRoomRes struct {
	ID               int                    `json:"id"`
	MeetingID        string                 `json:"meeting_id"`
	CreatorAccessKey string                 `json:"creator_access_key"`
	CreatorNickname  string                 `json:"creator_nickname,omitempty"` // 创建者昵称
	CallType         int8                   `json:"call_type"`
	CallTypeName     string                 `json:"call_type_name,omitempty"` // 通话类型名称
	Status           int8                   `json:"status"`
	StatusName       string                 `json:"status_name,omitempty"` // 状态名称
	Participants     map[string]interface{} `json:"participants"`
	ParticipantCount int                    `json:"participant_count,omitempty"` // 参与者数量
	CreatedAt        int64                  `json:"created_at"`
	UpdatedAt        int64                  `json:"updated_at"`
}

// MeetingParticipant 会议参与者结构
type MeetingParticipant struct {
	AccessKey  string `json:"access_key"`  // 访问密钥
	Nickname   string `json:"nickname"`    // 昵称
	JoinTime   int64  `json:"join_time"`   // 加入时间
	LeaveTime  int64  `json:"leave_time"`  // 离开时间
	Status     int8   `json:"status"`      // 状态：1在线 2离线
	DeviceType string `json:"device_type"` // 设备类型
}

// 常量定义

// 通话类型常量
const (
	MeetingCallTypeVoice int8 = 1 // 语音通话
	MeetingCallTypeVideo int8 = 2 // 视频通话
)

// 会议状态常量
const (
	MeetingStatusActive    int8 = 1 // 进行中
	MeetingStatusEnded     int8 = 2 // 已结束
	MeetingStatusCancelled int8 = 3 // 已取消
)

// 参与者状态常量
const (
	ParticipantStatusOnline  int8 = 1 // 在线
	ParticipantStatusOffline int8 = 2 // 离线
)

// GetMeetingCallTypeName 获取通话类型名称
func GetMeetingCallTypeName(callType int8) string {
	switch callType {
	case MeetingCallTypeVoice:
		return "语音通话"
	case MeetingCallTypeVideo:
		return "视频通话"
	default:
		return "未知类型"
	}
}

// GetMeetingStatusName 获取会议状态名称
func GetMeetingStatusName(status int8) string {
	switch status {
	case MeetingStatusActive:
		return "进行中"
	case MeetingStatusEnded:
		return "已结束"
	case MeetingStatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}

// GetParticipantStatusName 获取参与者状态名称
func GetParticipantStatusName(status int8) string {
	switch status {
	case ParticipantStatusOnline:
		return "在线"
	case ParticipantStatusOffline:
		return "离线"
	default:
		return "未知状态"
	}
}
