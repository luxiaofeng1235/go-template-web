package models

// SecretKey 用户密钥主模型
type SecretKey struct {
	ID                int64  `gorm:"column:id" json:"id"`                                 // ID
	DeviceFingerprint string `gorm:"column:device_fingerprint" json:"device_fingerprint"` // 设备指纹
	AccessKey         string `gorm:"column:access_key" json:"access_key"`                 // 访问密钥
	GroupID           string `gorm:"column:group_id" json:"group_id"`                     // 群组ID
	Nickname          string `gorm:"column:nickname" json:"nickname"`                     // 昵称
	UserNote          string `gorm:"column:user_note" json:"user_note"`                   // 用户备注
	AvtarURL          string `gorm:"column:avtar_url" json:"avtar_url"`                   // 头像URL
	Position          string `gorm:"column:position" json:"position"`                     // 位置信息
	DeviceInfo        string `gorm:"column:device_info" json:"device_info"`               // 设备信息
	FirstVisitTime    int64  `gorm:"column:first_visit_time" json:"first_visit_time"`     // 首次访问时间
	LastVisitTime     int64  `gorm:"column:last_visit_time" json:"last_visit_time"`       // 最后访问时间
	VisitCount        int64  `gorm:"column:visit_count" json:"visit_count"`               // 访问次数
	Status            int8   `gorm:"column:status" json:"status"`                         // 状态：1正常 2禁用 3删除
	CreatedAt         int64  `gorm:"column:created_at" json:"created_at"`                 // 创建时间
	UpdatedAt         int64  `gorm:"column:updated_at" json:"updated_at"`                 // 更新时间
}

func (*SecretKey) TableName() string {
	return "ls_secret_keys"
}

// 请求结构体

// CreateSecretKeyReq 创建用户密钥请求
type CreateSecretKeyReq struct {
	DeviceFingerprint string `form:"device_fingerprint" json:"device_fingerprint" `
	GroupID           string `form:"group_id" json:"group_id" `
	Nickname          string `form:"nickname" json:"nickname"`
	UserNote          string `form:"user_note" json:"user_note"`
	AvtarURL          string `form:"avtar_url" json:"avtar_url"`
	Position          string `form:"position" json:"position"`
	DeviceInfo        string `form:"device_info" json:"device_info"`
	CreatedAt         int64  `form:"created_at" json:"created_at"` // 创建时间
	UpdatedAt         int64  `form:"updated_at" json:"updated_at"` // 更新时间
}

// UpdateSecretKeyReq 更新用户密钥请求
type UpdateSecretKeyReq struct {
	AccessKey string `form:"access_key" json:"access_key" `
	CreateSecretKeyReq
}

// GetSecretKeyListReq 获取用户密钥列表请求
type GetSecretKeyListReq struct {
	GroupID  string `form:"group_id" json:"group_id"`
	Status   int8   `form:"status" json:"status"`
	Page     int    `form:"page" json:"page" `
	PageSize int    `form:"page_size" json:"page_size" `
	Keyword  string `form:"keyword" json:"keyword"` // 搜索关键词（昵称、备注）
}

// GetSecretKeyDetailReq 获取用户密钥详情请求
type GetSecretKeyDetailReq struct {
	AccessKey string `form:"access_key" json:"access_key" `
}

// UpdateVisitInfoReq 更新访问信息请求
type UpdateVisitInfoReq struct {
	AccessKey  string `form:"access_key" json:"access_key" `
	DeviceInfo string `form:"device_info" json:"device_info"`
}

// 响应结构体

// SecretKeyRes 用户密钥响应
type SecretKeyRes struct {
	ID                int64  `json:"id"`
	DeviceFingerprint string `json:"device_fingerprint"`
	AccessKey         string `json:"access_key"`
	GroupID           string `json:"group_id"`
	Nickname          string `json:"nickname"`
	UserNote          string `json:"user_note"`
	AvtarURL          string `json:"avtar_url"`
	Position          string `json:"position"`
	FirstVisitTime    int64  `json:"first_visit_time"`
	LastVisitTime     int64  `json:"last_visit_time"`
	VisitCount        int64  `json:"visit_count"`
	Status            int8   `json:"status"`
	StatusName        string `json:"status_name"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
	IsOnline          bool   `json:"is_online"` // 是否在线（需要结合session判断）
}

// SecretKeyListRes 用户密钥列表响应
type SecretKeyListRes struct {
	SecretKeys []SecretKeyRes `json:"secret_keys"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	HasMore    bool           `json:"has_more"`
}

// CreateSecretKeyRes 创建用户密钥响应
type CreateSecretKeyRes struct {
	AccessKey string `json:"access_key"`
	GroupID   string `json:"group_id"`
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
}

// GroupStatsRes 群组统计响应
type GroupStatsRes struct {
	GroupID     string `json:"group_id"`
	TotalUsers  int64  `json:"total_users"`
	OnlineUsers int64  `json:"online_users"`
	ActiveUsers int64  `json:"active_users"` // 今日活跃用户
}

// 常量定义

// 用户状态常量
const (
	SecretKeyStatusNormal   int8 = 1 // 正常
	SecretKeyStatusDisabled int8 = 2 // 禁用
	SecretKeyStatusDeleted  int8 = 3 // 删除
)

// 获取状态名称
func GetSecretKeyStatusName(status int8) string {
	switch status {
	case SecretKeyStatusNormal:
		return "正常"
	case SecretKeyStatusDisabled:
		return "禁用"
	case SecretKeyStatusDeleted:
		return "删除"
	default:
		return "未知状态"
	}
}
