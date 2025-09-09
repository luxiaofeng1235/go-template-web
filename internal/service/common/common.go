package common

import (
	"go-web-template/internal/constant"
)

// GetProductCategoryList 获取商品分类列表
func GetProductCategoryList() []constant.ProductCategory {
	// 返回从constant中定义的分类列表
	return constant.ProductCategoryList
}

// GetProductCategoryByValue 根据值获取商品分类信息
func GetProductCategoryByValue(value int) *constant.ProductCategory {
	for _, category := range constant.ProductCategoryList {
		if category.Value == value {
			return &category
		}
	}
	return nil
}

// GetProductCategoryLabel 根据值获取分类名称
func GetProductCategoryLabel(value int) string {
	category := GetProductCategoryByValue(value)
	if category != nil {
		return category.Label
	}
	return "未知分类"
}