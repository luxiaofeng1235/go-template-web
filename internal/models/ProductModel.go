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

// TableName 指定表名
func (Product) TableName() string {
	return "ls_product"
}