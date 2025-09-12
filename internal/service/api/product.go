/*
 * @file: product.go
 * @description: 商品相关业务逻辑处理 - API层服务
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package api

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"

	"github.com/gogf/gf/v2/net/ghttp"
)

// GetProductList 获取商品列表
// @param r *ghttp.Request HTTP请求对象，用于获取域名信息
// @param req *models.ProductListReq 商品列表查询请求参数
// @return list []models.ProductListItem 商品列表数据
// @return err error 错误信息
func GetProductList(r *ghttp.Request, req *models.ProductListReq) (list []models.ProductListItem, err error) {
	// 参数验证
	if req == nil {
		err = fmt.Errorf("请求参数不能为空")
		global.Errlog.Error("商品列表查询参数无效", "req", req)
		return
	}
	// 参数默认值处理
	if req.PageNo <= 0 {
		req.PageNo = constant.PAGE_NO
	}
	if req.PageSize <= 0 {
		req.PageSize = constant.PAGE_SIZE
	}
	global.Requestlog.Info("开始查询商品列表",
		"pageNo", req.PageNo,
		"pageSize", req.PageSize,
		"cateID", req.CateID)
	// 构建查询条件
	query := global.DB.Model(&models.Product{}).
		Select("id, product_name, cate_id, intro, logo, qrcode")
	// 添加分类筛选条件（如果提供了cate_id）
	if req.CateID > 0 {
		query = query.Where("cate_id = ?", req.CateID)
		global.Sqllog.Info("添加分类筛选条件", "cateID", req.CateID)
	}
	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	var products []models.Product

	err = query.Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&products).Error

	if err != nil {
		global.Errlog.Error("查询商品列表失败",
			"pageNo", req.PageNo,
			"pageSize", req.PageSize,
			"cateID", req.CateID,
			"error", err)
		return
	}

	// 检查查询结果
	if len(products) == 0 {
		global.Requestlog.Info("商品列表查询结果为空",
			"pageNo", req.PageNo,
			"cateID", req.CateID)
		return []models.ProductListItem{}, nil
	}

	// 转换为返回格式并处理数据
	list = make([]models.ProductListItem, 0, len(products))

	for _, product := range products {
		// 检查商品数据有效性
		if product.ID <= 0 {
			continue
		}
		item := models.ProductListItem{
			ID:          product.ID,
			ProductName: product.ProductName,
			CateID:      product.CateID,
			Intro:       product.Intro,
			Logo:        utils.ProcessLogoURLForStatic(product.Logo, r),
			QRCode:      product.QRCode,
			CateName:    constant.GetProductCategoryName(product.CateID),
		}
		list = append(list, item)
	}

	global.Requestlog.Info("商品列表查询成功",
		"pageNo", req.PageNo,
		"pageSize", req.PageSize,
		"cateID", req.CateID,
		"count", len(list))

	return list, nil
}
