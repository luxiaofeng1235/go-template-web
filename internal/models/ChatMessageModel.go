package models

// ChatMessage 聊天消息主模型
type ChatMessage struct {
	ID          int64  `gorm:"column:id" json:"id"`                     // ID
	MessageID   string `gorm:"column:message_id" json:"message_id"`     // 消息ID
	SecretKey   string `gorm:"column:secret_key" json:"secret_key"`     // 密钥
	SenderID    string `gorm:"column:sender_id" json:"sender_id"`       // 发送者ID
	SenderName  string `gorm:"column:sender_name" json:"sender_name"`   // 发送者名称
	ReceiverID  string `gorm:"column:receiver_id" json:"receiver_id"`   // 接收者ID
	MessageType int8   `gorm:"column:message_type" json:"message_type"` // 消息类型：1文本 2图片 3视频 4文件
	Content     string `gorm:"column:content" json:"content"`           // 内容
	ImageURL    string `gorm:"column:image_url" json:"image_url"`       // 图片URL
	VideoURL    string `gorm:"column:video_url" json:"video_url"`       // 视频URL
	ChatType    int8   `gorm:"column:chat_type" json:"chat_type"`       // 聊天类型：0私聊 1群聊 2系统消息
	ExtraData   string `gorm:"column:extra_data" json:"extra_data"`     // 额外数据
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   int64  `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

func (*ChatMessage) TableName() string {
	return "ls_chat_messages"
}

// 请求结构体

// SendMessageReq 发送消息请求
type SendMessageReq struct {
	SecretKey   string `form:"secret_key" json:"secret_key" `
	ReceiverID  string `form:"receiver_id" json:"receiver_id"`
	MessageType int8   `form:"message_type" json:"message_type" `
	Content     string `form:"content" json:"content"`
	ImageURL    string `form:"image_url" json:"image_url"`
	VideoURL    string `form:"video_url" json:"video_url"`
	ChatType    int8   `form:"chat_type" json:"chat_type" `
	ExtraData   string `form:"extra_data" json:"extra_data"`
}

// GetMessageListReq 获取消息列表请求
type GetMessageListReq struct {
	SecretKey  string `form:"secret_key" json:"secret_key" `
	ChatType   int8   `form:"chat_type" json:"chat_type"`
	ReceiverID string `form:"receiver_id" json:"receiver_id"`
	Page       int    `form:"page" json:"page" `
	PageSize   int    `form:"page_size" json:"page_size" `
	LastMsgID  string `form:"last_msg_id" json:"last_msg_id"`
}

// GetMessageDetailReq 获取消息详情请求
type GetMessageDetailReq struct {
	MessageID string `form:"message_id" json:"message_id" `
	SecretKey string `form:"secret_key" json:"secret_key" `
}

// ChatUserListReq 获取聊天用户列表请求（对应chat/getDeviceList）
type ChatUserListReq struct {
	AccessKey string `form:"access_key" json:"access_key" v:"required#访问密钥不能为空"`
	PageNo    int    `form:"page_no" json:"page_no"`
	PageSize  int    `form:"page_size" json:"page_size"`
}

// ChatHistoryReq 获取聊天历史请求（对应Chat/getChatHistory）
type ChatHistoryReq struct {
	SecretKey   string `form:"secret_key" json:"secret_key" v:"required#密钥不能为空"`
	UserId      int64  `form:"user_id" json:"user_id" v:"required#用户ID不能为空"`
	ToUserId    int64  `form:"to_user_id" json:"to_user_id" v:"required#对方用户ID不能为空"`
	PageNum     int    `form:"page_num" json:"page_num"`
	PageSize    int    `form:"page_size" json:"page_size"`
	StartTime   int64  `form:"start_time" json:"start_time"`
	EndTime     int64  `form:"end_time" json:"end_time"`
	MessageType int8   `form:"message_type" json:"message_type"`
}

// 响应结构体

// ChatMessageRes 聊天消息响应
type ChatMessageRes struct {
	ID          int64  `json:"id"`
	MessageID   string `json:"message_id"`
	SenderID    string `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	ReceiverID  string `json:"receiver_id"`
	MessageType int8   `json:"message_type"`
	Content     string `json:"content"`
	ImageURL    string `json:"image_url,omitempty"`
	VideoURL    string `json:"video_url,omitempty"`
	ChatType    int8   `json:"chat_type"`
	ExtraData   string `json:"extra_data,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// MessageListRes 消息列表响应
type MessageListRes struct {
	Messages  []ChatMessageRes `json:"messages"`
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	HasMore   bool             `json:"has_more"`
	LastMsgID string           `json:"last_msg_id"`
}

// SendMessageRes 发送消息响应
type SendMessageRes struct {
	MessageID   int64  `json:"message_id"`
	UserId      int64  `json:"user_id"`
	ToUserId    int64  `json:"to_user_id"`
	Content     string `json:"content"`
	MessageType int8   `json:"message_type"`
	CreateTime  int64  `json:"create_time"`
}

// ChatUserListRes 聊天用户列表响应（对应chat/getDeviceList）
type ChatUserListRes struct {
	DeviceId        string `json:"device_id"`
	UserId          int64  `json:"user_id"`
	Username        string `json:"username"`
	LastMessageTime int64  `json:"last_message_time"`
	MessageCount    int64  `json:"message_count"`
	LastMessage     string `json:"last_message"`
	UnreadCount     int64  `json:"unread_count"`
}

// ChatHistoryRes 聊天历史响应（对应Chat/getChatHistory）
type ChatHistoryRes struct {
	Id          int64  `json:"id"`
	MessageId   string `json:"message_id"`
	UserId      int64  `json:"user_id"`
	ToUserId    int64  `json:"to_user_id"`
	Username    string `json:"username"`
	ToUsername  string `json:"to_username"`
	Content     string `json:"content"`
	MessageType int8   `json:"message_type"`
	ImageUrl    string `json:"image_url,omitempty"`
	VideoUrl    string `json:"video_url,omitempty"`
	ChatType    int8   `json:"chat_type"`
	Status      int8   `json:"status"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
}
