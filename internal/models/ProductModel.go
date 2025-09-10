package models

// Product 商品表结构体
type Product struct {
	ID          int64  `gorm:"column:id" json:"id"`                     // 商品ID
	ProductName string `gorm:"column:product_name" json:"product_name"` // 商品名称
	CateID      int    `gorm:"column:cate_id" json:"cate_id"`           // 分类ID
	Intro       string `gorm:"column:intro" json:"intro"`               // 商品介绍
	Logo        string `gorm:"column:logo" json:"logo"`                 // 商品Logo
	QRCode      string `gorm:"column:qrcode" json:"qrcode"`             // 二维码
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   int64  `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

// ProductListReq 商品列表请求参数
type ProductListReq struct {
	PageNo   int `json:"page_no" form:"page_no" v:"required|min:1#页码必须输入|页码必须大于0"`         // 页码
	PageSize int `json:"page_size" form:"page_size" v:"required|between:1,100#每页数量必须输入|每页数量必须在1-100之间"` // 每页数量
	CateID   int `json:"cate_id" form:"cate_id" v:"min:0#分类ID不能小于0"`                      // 分类ID（可选）
}

// ProductListItem 商品列表项
type ProductListItem struct {
	ID          int64  `json:"id"`           // 商品ID
	ProductName string `json:"product_name"` // 商品名称
	CateID      int    `json:"cate_id"`      // 分类ID
	CateName    string `json:"cate_name"`    // 分类名称
	Intro       string `json:"intro"`        // 商品介绍
	Logo        string `json:"logo"`         // 商品Logo（完整URL）
	QRCode      string `json:"qrcode"`       // 二维码
}

// SaveProductReq 保存商品请求参数
type SaveProductReq struct {
	ID          int64  `json:"id" form:"id" v:"min:0#商品ID不能小于0"`                                     // 商品ID（编辑时传入，新增时为空）
	ProductName string `json:"product_name" form:"product_name" v:"required|max-length:100#商品名称必须输入|商品名称长度不能超过100字符"` // 商品名称（必填）
	CateID      int    `json:"cate_id" form:"cate_id" v:"required|min:1#分类ID必须选择|分类ID必须大于0"`           // 分类ID
	Intro       string `json:"intro" form:"intro" v:"max-length:500#商品介绍长度不能超过500字符"`                   // 商品介绍
	Logo        string `json:"logo" form:"logo" v:"max-length:255#商品Logo链接长度不能超过255字符"`                 // 商品Logo
	QRCode      string `json:"qrcode" form:"qrcode" v:"max-length:255#二维码链接长度不能超过255字符"`              // 二维码
}

// SaveProductRes 保存商品响应数据
type SaveProductRes struct {
	ID          int64  `json:"id"`           // 商品ID
	ProductName string `json:"product_name"` // 商品名称
	CateID      int    `json:"cate_id"`      // 分类ID
	CreatedAt   int64  `json:"created_at"`   // 创建时间
	UpdatedAt   int64  `json:"updated_at"`   // 更新时间
}

// TableName 指定表名
func (Product) TableName() string {
	return "ls_product"
}
