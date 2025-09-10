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
	Status            int    `gorm:"column:status" json:"status"`                         // 状态：1正常 2禁用 3删除
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
	DeviceInfo        string `form:"device_info" json:"device_info"`
}

// UpdateSecretKeyReq 更新用户密钥请求
type UpdateSecretKeyReq struct {
	AccessKey string `form:"access_key" json:"access_key" `
	CreateSecretKeyReq
}

// 创建密钥的返回值
type CreateSecretResp struct {
	AccessKey string `form:"access_key" json:"access_key" `
	GROUPID   string `form:"group_id" json:"group_id" `
	AVTAR_URL string `form:"avtar_url" json:"avtar_url" `
	NickName  string `form:"nick_name" json:"nick_name" `
	IsNew     bool   `form:"user_note" json:"user_note" `
}
