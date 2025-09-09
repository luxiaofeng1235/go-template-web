package models

import (
	"time"
)

// ProductModel 产品数据模型
type ProductModel struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null;comment:产品名称"`
	Description string    `json:"description" gorm:"type:text;comment:产品描述"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);default:0.00;comment:价格"`
	Stock       int       `json:"stock" gorm:"type:int;default:0;comment:库存"`
	Image       string    `json:"image" gorm:"type:varchar(255);comment:产品图片"`
	CategoryID  int       `json:"category_id" gorm:"type:int;default:0;comment:分类ID"`
	Status      int       `json:"status" gorm:"type:tinyint;default:1;comment:状态:1正常,0禁用"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (ProductModel) TableName() string {
	return "ls_product"
}

// ProductCreateReq 创建产品请求
type ProductCreateReq struct {
	Name        string  `json:"name" v:"required|length:1,100#产品名称不能为空|产品名称长度为1-100位"`
	Description string  `json:"description"`
	Price       float64 `json:"price" v:"required|min:0#价格不能为空|价格不能小于0"`
	Stock       int     `json:"stock" v:"min:0#库存不能小于0"`
	Image       string  `json:"image"`
	CategoryID  int     `json:"category_id"`
}

// ProductUpdateReq 更新产品请求
type ProductUpdateReq struct {
	ID          int     `json:"id" v:"required|min:1#产品ID不能为空|产品ID无效"`
	Name        string  `json:"name" v:"required|length:1,100#产品名称不能为空|产品名称长度为1-100位"`
	Description string  `json:"description"`
	Price       float64 `json:"price" v:"required|min:0#价格不能为空|价格不能小于0"`
	Stock       int     `json:"stock" v:"min:0#库存不能小于0"`
	Image       string  `json:"image"`
	CategoryID  int     `json:"category_id"`
}

// ProductListRes 产品列表响应
type ProductListRes struct {
	List  []ProductModel `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
}
