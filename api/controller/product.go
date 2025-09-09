package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"go-web-template/internal/service"
	"go-web-template/utils"
)

type ProductController struct{}

// GetList 获取产品列表
func (c *ProductController) GetList(r *ghttp.Request) {
	page := r.Get("page", 1).Int()
	pageSize := r.Get("page_size", 10).Int()

	products, total, err := service.Product.GetList(r.Context(), page, pageSize)
	if err != nil {
		utils.FailEncrypt(r, err, "获取产品列表失败")
		return
	}

	utils.SuccessPage(r, products, total, page, pageSize, "获取成功")
}

// GetDetail 获取产品详情
func (c *ProductController) GetDetail(r *ghttp.Request) {
	id := r.Get("id").Int()
	if id <= 0 {
		utils.ParamError(r, "产品ID无效")
		return
	}

	product, err := service.Product.GetDetail(r.Context(), id)
	if err != nil {
		utils.FailEncrypt(r, err, "获取产品详情失败")
		return
	}

	utils.SuccessEncrypt(r, product, "获取成功")
}

// Create 创建产品
func (c *ProductController) Create(r *ghttp.Request) {
	name := r.Get("name").String()
	description := r.Get("description").String()
	price := r.Get("price").Float64()

	if name == "" {
		utils.ParamError(r, "产品名称不能为空")
		return
	}

	err := service.Product.Create(r.Context(), name, description, price)
	if err != nil {
		utils.Fail(r, err, "创建产品失败")
		return
	}

	utils.Success(r, nil, "创建成功")
}