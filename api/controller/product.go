package controller

import (
	"go-web-template/utils"

	"github.com/gogf/gf/v2/net/ghttp"
)

type ProductController struct{}

// GetProductList 获取商品列表
func (c *ProductController) GetProductList(r *ghttp.Request) {
	// TODO: 实现商品列表逻辑
	utils.Success(r, map[string]interface{}{
		"list":     []interface{}{},
		"total":    0,
		"page":     1,
		"pageSize": 10,
	}, "获取商品列表成功")
}

// GetCategoryList 获取分类列表
func (c *ProductController) GetCategoryList(r *ghttp.Request) {
	// TODO: 实现分类列表逻辑
	utils.Success(r, map[string]interface{}{
		"categories": []interface{}{},
		"total":      0,
	}, "获取分类列表成功")
}