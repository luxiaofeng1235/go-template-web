package models

// File 文件主模型
type File struct {
	ID         uint   `gorm:"column:id" json:"id"`                   // ID
	Name       string `gorm:"column:name" json:"name"`               // 文件名
	CID        uint   `gorm:"column:cid" json:"cid"`                 // 分类ID
	Type       uint8  `gorm:"column:type" json:"type"`               // 文件类型
	URI        string `gorm:"column:uri" json:"uri"`                 // 文件路径
	CreateTime uint   `gorm:"column:create_time" json:"create_time"` // 创建时间
	Del        uint8  `gorm:"column:del" json:"del"`                 // 删除状态：0正常 1删除
	ShopID     int    `gorm:"column:shop_id" json:"shop_id"`         // 店铺ID
	UserID     uint   `gorm:"column:user_id" json:"user_id"`         // 用户ID
}

func (*File) TableName() string {
	return "ls_file"
}

// 请求结构体

// CreateFileReq 创建文件请求
type CreateFileReq struct {
	Name   string `form:"name" json:"name"`
	CID    uint   `form:"cid" json:"cid"`
	Type   uint8  `form:"type" json:"type"`
	URI    string `form:"uri" json:"uri"`
	ShopID int    `form:"shop_id" json:"shop_id"`
	UserID uint   `form:"user_id" json:"user_id"`
}

// UpdateFileReq 更新文件请求
type UpdateFileReq struct {
	ID uint `form:"id" json:"id"`
	CreateFileReq
}

// GetFileListReq 获取文件列表请求
type GetFileListReq struct {
	CID      uint   `form:"cid" json:"cid"`
	Type     uint8  `form:"type" json:"type"`
	Name     string `form:"name" json:"name"`
	ShopID   int    `form:"shop_id" json:"shop_id"`
	UserID   uint   `form:"user_id" json:"user_id"`
	Del      uint8  `form:"del" json:"del"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

// GetFileDetailReq 获取文件详情请求
type GetFileDetailReq struct {
	ID uint `form:"id" json:"id"`
}

// 响应结构体

// FileRes 文件响应
type FileRes struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	CID          uint   `json:"cid"`
	CategoryName string `json:"category_name,omitempty"` // 分类名称
	Type         uint8  `json:"type"`
	TypeName     string `json:"type_name,omitempty"` // 类型名称
	URI          string `json:"uri"`
	CreateTime   uint   `json:"create_time"`
	Del          uint8  `json:"del"`
	DelName      string `json:"del_name,omitempty"` // 删除状态名称
	ShopID       int    `json:"shop_id"`
	UserID       uint   `json:"user_id"`
}

// 常量定义

// 文件类型常量
const (
	FileTypeImage    uint8 = 10 // 图片
	FileTypeDocument uint8 = 20 // 文档
	FileTypeVideo    uint8 = 30 // 视频
	FileTypeAudio    uint8 = 40 // 音频
	FileTypeArchive  uint8 = 50 // 压缩包
	FileTypeOther    uint8 = 60 // 其他
)

// 删除状态常量
const (
	FileDelNormal  uint8 = 0 // 正常
	FileDelDeleted uint8 = 1 // 删除
)

// GetFileTypeName 获取文件类型名称
func GetFileTypeName(fileType uint8) string {
	switch fileType {
	case FileTypeImage:
		return "图片"
	case FileTypeDocument:
		return "文档"
	case FileTypeVideo:
		return "视频"
	case FileTypeAudio:
		return "音频"
	case FileTypeArchive:
		return "压缩包"
	case FileTypeOther:
		return "其他"
	default:
		return "未知类型"
	}
}

// GetFileDelName 获取删除状态名称
func GetFileDelName(del uint8) string {
	switch del {
	case FileDelNormal:
		return "正常"
	case FileDelDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}
