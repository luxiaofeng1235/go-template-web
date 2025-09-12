/*
 * @file: product.go
 * @description: 商品管理控制器 - Admin后台管理模块
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package controller

import (
	"go-web-template/internal/models"
	"go-web-template/internal/service/admin"
	"go-web-template/utils"
	"math"

	"github.com/gogf/gf/v2/net/ghttp"
)

type ProductController struct{}

// GetProductList 获取商品列表 - 后台管理
func (c *ProductController) GetProductList(r *ghttp.Request) {
	// 解析请求参数
	var req models.ProductListReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析错误")
		return
	}

	// 调用Service层处理业务逻辑
	list, total, err := admin.GetProductList(r, &req)
	if err != nil {
		utils.FailEncrypt(r, err, "获取商品列表失败")
		return
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	// 构造响应数据 - 按照PHP版本格式
	data := map[string]interface{}{
		"total":       total,
		"list":        list,
		"page_no":     req.PageNo,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	utils.Success(r, data, "获取商品列表成功")
}

// GetCategoryList 获取分类列表 - 后台管理
func (c *ProductController) GetCategoryList(r *ghttp.Request) {
	// 调用Service层处理业务逻辑
	categories, err := admin.GetCategoryList()
	if err != nil {
		utils.FailEncrypt(r, err, "获取分类列表失败")
		return
	}

	// 直接返回分类数组，不需要包装在categories字段中
	utils.Success(r, categories, "分类列表")
}

// SaveProduct 保存商品信息 - 支持新增和编辑
func (c *ProductController) SaveProduct(r *ghttp.Request) {
	// 解析请求参数
	var req models.SaveProductReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析错误")
		return
	}

	// 调用Service层处理业务逻辑
	result, err := admin.SaveProduct(&req)
	if err != nil {
		utils.FailEncrypt(r, err, "保存商品失败")
		return
	}

	utils.Success(r, result, "保存商品成功")
}

// DeleteProduct 删除商品
func (c *ProductController) DeleteProduct(r *ghttp.Request) {
	// 获取商品ID参数
	productID := r.Get("id").Int64()
	if productID <= 0 {
		utils.FailEncrypt(r, nil, "商品ID不能为空")
		return
	}

	// 调用Service层处理业务逻辑
	err := admin.DeleteProduct(productID)
	if err != nil {
		utils.FailEncrypt(r, err, "删除商品失败")
		return
	}

	utils.Success(r, nil, "删除商品成功")
}
