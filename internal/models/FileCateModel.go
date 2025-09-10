package models

import "go-web-template/internal/constant"

// FileCate 文件分类主模型
type FileCate struct {
	ID         uint   `gorm:"column:id" json:"id"`                   // ID
	ShopID     int    `gorm:"column:shop_id" json:"shop_id"`         // 店铺ID
	Name       string `gorm:"column:name" json:"name"`               // 分类名称
	PID        uint   `gorm:"column:pid" json:"pid"`                 // 父级ID
	Type       uint8  `gorm:"column:type" json:"type"`               // 分类类型
	Level      int8   `gorm:"column:level" json:"level"`             // 层级
	Sort       uint16 `gorm:"column:sort" json:"sort"`               // 排序
	Del        uint8  `gorm:"column:del" json:"del"`                 // 删除状态：0正常 1删除
	CreateTime int64  `gorm:"column:create_time" json:"create_time"` // 创建时间
	UpdateTime int64  `gorm:"column:update_time" json:"update_time"` // 更新时间
}

func (*FileCate) TableName() string {
	return "ls_file_cate"
}

// 请求结构体

// CreateFileCateReq 创建文件分类请求
type CreateFileCateReq struct {
	ShopID     int    `form:"shop_id" json:"shop_id"`
	Name       string `form:"name" json:"name"`
	PID        uint   `form:"pid" json:"pid"`
	Type       uint8  `form:"type" json:"type"`
	Level      int8   `form:"level" json:"level"`
	Sort       uint16 `form:"sort" json:"sort"`
	CreateTime int64  `form:"create_time" json:"create_time"` // 创建时间
	UpdateTime int64  `form:"update_time" json:"update_time"` // 更新时间
}

// UpdateFileCateReq 更新文件分类请求
type UpdateFileCateReq struct {
	ID uint `form:"id" json:"id"`
	CreateFileCateReq
}

// GetFileCateListReq 获取文件分类列表请求
type GetFileCateListReq struct {
	ShopID   int    `form:"shop_id" json:"shop_id"`
	Name     string `form:"name" json:"name"`
	PID      uint   `form:"pid" json:"pid"`
	Type     uint8  `form:"type" json:"type"`
	Level    int8   `form:"level" json:"level"`
	Del      uint8  `form:"del" json:"del"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// GetFileCateDetailReq 获取文件分类详情请求
type GetFileCateDetailReq struct {
	ID uint `form:"id" json:"id"`
}

// 响应结构体

// FileCateRes 文件分类响应
type FileCateRes struct {
	ID         uint          `json:"id"`
	ShopID     int           `json:"shop_id"`
	Name       string        `json:"name"`
	PID        uint          `json:"pid"`
	ParentName string        `json:"parent_name,omitempty"` // 父级名称
	Type       uint8         `json:"type"`
	TypeName   string        `json:"type_name,omitempty"` // 类型名称
	Level      int8          `json:"level"`
	Sort       uint16        `json:"sort"`
	Del        uint8         `json:"del"`
	DelName    string        `json:"del_name,omitempty"` // 删除状态名称
	CreateTime int64         `json:"create_time"`
	UpdateTime int64         `json:"update_time"`
	Children   []FileCateRes `json:"children,omitempty"` // 子分类
}

// GetFileCateTypeName 获取文件分类类型名称
func GetFileCateTypeName(cateType uint8) string {
	switch cateType {
	case constant.FileCateTypeImage:
		return "图片分类"
	case constant.FileCateTypeDocument:
		return "文档分类"
	case constant.FileCateTypeVideo:
		return "视频分类"
	case constant.FileCateTypeAudio:
		return "音频分类"
	case constant.FileCateTypeArchive:
		return "压缩包分类"
	case constant.FileCateTypeOther:
		return "其他分类"
	default:
		return "未知类型"
	}
}

// GetFileCateDelName 获取文件分类删除状态名称
func GetFileCateDelName(del uint8) string {
	switch del {
	case constant.FileCateDelNormal:
		return "正常"
	case constant.FileCateDelDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}
