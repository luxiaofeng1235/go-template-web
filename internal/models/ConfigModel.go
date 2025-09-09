package models

// Config 配置主模型
type Config struct {
	ID         int    `gorm:"column:id" json:"id"`                   // ID
	ShopID     int    `gorm:"column:shop_id" json:"shop_id"`         // 店铺ID
	Type       string `gorm:"column:type" json:"type"`               // 配置类型
	Name       string `gorm:"column:name" json:"name"`               // 配置名称
	Value      string `gorm:"column:value" json:"value"`             // 配置值
	CreatedAt  int64  `gorm:"column:created_at" json:"created_at"`   // 创建时间
	UpdateTime int64  `gorm:"column:update_time" json:"update_time"` // 更新时间
}

func (*Config) TableName() string {
	return "ls_config"
}

// 请求结构体

// CreateConfigReq 创建配置请求
type CreateConfigReq struct {
	ShopID     int    `form:"shop_id" json:"shop_id"`
	Type       string `form:"type" json:"type"`
	Name       string `form:"name" json:"name"`
	Value      string `form:"value" json:"value"`
	CreatedAt  int64  `form:"created_at" json:"created_at"`   // 创建时间
	UpdateTime int64  `form:"update_time" json:"update_time"` // 更新时间
}

// UpdateConfigReq 更新配置请求
type UpdateConfigReq struct {
	ID int `form:"id" json:"id"`
	CreateConfigReq
}

// GetConfigListReq 获取配置列表请求
type GetConfigListReq struct {
	ShopID   int    `form:"shop_id" json:"shop_id"`
	Type     string `form:"type" json:"type"`
	Name     string `form:"name" json:"name"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// GetConfigDetailReq 获取配置详情请求
type GetConfigDetailReq struct {
	ID int `form:"id" json:"id"`
}

// 响应结构体

// ConfigRes 配置响应
type ConfigRes struct {
	ID         int    `json:"id"`
	ShopID     int    `json:"shop_id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	CreatedAt  int64  `json:"created_at"`
	UpdateTime int64  `json:"update_time"`
}

// 常量定义

// 配置类型常量
const (
	ConfigTypeSystem   = "system"   // 系统配置
	ConfigTypeUser     = "user"     // 用户配置
	ConfigTypeShop     = "shop"     // 店铺配置
	ConfigTypeFeature  = "feature"  // 功能配置
	ConfigTypePayment  = "payment"  // 支付配置
	ConfigTypeSecurity = "security" // 安全配置
)

// GetConfigTypeName 获取配置类型名称
func GetConfigTypeName(configType string) string {
	switch configType {
	case ConfigTypeSystem:
		return "系统配置"
	case ConfigTypeUser:
		return "用户配置"
	case ConfigTypeShop:
		return "店铺配置"
	case ConfigTypeFeature:
		return "功能配置"
	case ConfigTypePayment:
		return "支付配置"
	case ConfigTypeSecurity:
		return "安全配置"
	default:
		return "未知类型"
	}
}
