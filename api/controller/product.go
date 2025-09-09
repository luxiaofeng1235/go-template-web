package controller

import (
	"go-web-template/internal/service/common"
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
	// 调用common服务获取分类列表
	categories := common.GetProductCategoryList()
	
	// 将value=0的标签强制设置为"全部"
	for i, category := range categories {
		if category.Value == 0 {
			categories[i].Label = "全部"
			break
		}
	}
	
	// 直接返回分类数组作为data
	utils.Success(r, categories, "获取分类列表成功")
}
