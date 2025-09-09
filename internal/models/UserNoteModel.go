package models

// UserNote 用户备注主模型
type UserNote struct {
	ID          int    `gorm:"column:id" json:"id"`                       // ID
	AccessKey   string `gorm:"column:access_key" json:"access_key"`       // 访问密钥
	ToAccessKey string `gorm:"column:to_access_key" json:"to_access_key"` // 目标用户访问密钥
	UserNote    string `gorm:"column:user_note" json:"user_note"`         // 用户备注
	CreatedAt   int    `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt   int    `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
}

func (*UserNote) TableName() string {
	return "ls_user_note"
}

// 请求结构体

// CreateUserNoteReq 创建用户备注请求
type CreateUserNoteReq struct {
	AccessKey   string `form:"access_key" json:"access_key"`
	ToAccessKey string `form:"to_access_key" json:"to_access_key"`
	UserNote    string `form:"user_note" json:"user_note"`
}

// UpdateUserNoteReq 更新用户备注请求
type UpdateUserNoteReq struct {
	ID int `form:"id" json:"id"`
	CreateUserNoteReq
}

// GetUserNoteListReq 获取用户备注列表请求
type GetUserNoteListReq struct {
	AccessKey   string `form:"access_key" json:"access_key"`
	ToAccessKey string `form:"to_access_key" json:"to_access_key"`
	UserNote    string `form:"user_note" json:"user_note"`
	Page        int    `form:"page" json:"page"`
	PageSize    int    `form:"page_size" json:"page_size"`
}

// GetUserNoteDetailReq 获取用户备注详情请求
type GetUserNoteDetailReq struct {
	ID          int    `form:"id" json:"id"`
	AccessKey   string `form:"access_key" json:"access_key"`
	ToAccessKey string `form:"to_access_key" json:"to_access_key"`
}

// GetUserNoteByUsersReq 根据用户关系获取备注请求
type GetUserNoteByUsersReq struct {
	AccessKey   string `form:"access_key" json:"access_key"`
	ToAccessKey string `form:"to_access_key" json:"to_access_key"`
}

// 响应结构体

// UserNoteRes 用户备注响应
type UserNoteRes struct {
	ID             int    `json:"id"`
	AccessKey      string `json:"access_key"`
	ToAccessKey    string `json:"to_access_key"`
	UserNote       string `json:"user_note"`
	CreatedAt      int    `json:"created_at"`
	UpdatedAt      int    `json:"updated_at"`
	UserNickname   string `json:"user_nickname,omitempty"`    // 用户昵称
	ToUserNickname string `json:"to_user_nickname,omitempty"` // 目标用户昵称
}
