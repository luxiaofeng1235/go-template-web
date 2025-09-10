package controller

import (
	"go-web-template/internal/models"
	"go-web-template/internal/service/api"
	"go-web-template/internal/service/common"
	"go-web-template/utils"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ProductController struct{}

// GetProductList 获取商品列表
func (c *ProductController) GetProductList(r *ghttp.Request) {
	var req *models.ProductListReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 调用service层处理业务逻辑
	list, err := api.GetProductList(r, req)
	if err != nil {
		utils.Fail(r, err, "获取商品列表失败")
		return
	}

	utils.Success(r, list, "获取商品列表成功")
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
