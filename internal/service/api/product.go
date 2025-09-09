package api

import (
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ProductListReq 商品列表请求参数
type ProductListReq struct {
	PageNo   int `json:"page_no" form:"page_no"`     // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页数量
	CateID   int `json:"cate_id" form:"cate_id"`     // 分类ID（可选）
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

// GetProductList 获取商品列表
func GetProductList(r *ghttp.Request, req *ProductListReq) ([]ProductListItem, error) {
	// 参数默认值处理
	if req.PageNo <= 0 {
		req.PageNo = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	query := global.DB.Model(&models.Product{}).
		Select("id, product_name, cate_id, intro, logo, qrcode")

	// 添加分类筛选条件（如果提供了cate_id）
	if req.CateID > 0 {
		query = query.Where("cate_id = ?", req.CateID)
	}

	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	var products []models.Product
	err := query.Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&products).Error

	if err != nil {
		global.Errlog.Error("查询商品列表失败", "error", err)
		return nil, err
	}

	// 转换为返回格式并处理数据
	var result []ProductListItem
	baseURL := utils.GetFullDomain(r) // 获取当前域名

	for _, product := range products {
		item := ProductListItem{
			ID:          product.ID,
			ProductName: product.ProductName,
			CateID:      int(product.CateID),
			Intro:       product.Intro,
			Logo:        processLogoURL(product.Logo, baseURL),
			QRCode:      product.QRCode,
			CateName:    getCategoryName(int(product.CateID)),
		}
		result = append(result, item)
	}

	global.Requestlog.Info("商品列表查询成功", 
		"pageNo", req.PageNo, 
		"pageSize", req.PageSize, 
		"cateID", req.CateID,
		"count", len(result))

	return result, nil
}

// processLogoURL 处理Logo URL，确保返回完整URL
func processLogoURL(logo, baseURL string) string {
	if logo == "" {
		return ""
	}

	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(logo, "http://") || strings.HasPrefix(logo, "https://") {
		return logo
	}

	// 如果是相对路径，拼接域名
	if strings.HasPrefix(logo, "/") {
		return baseURL + logo
	}

	// 其他情况，添加/前缀后拼接
	return baseURL + "/" + logo
}

// getCategoryName 根据分类ID获取分类名称
func getCategoryName(cateID int) string {
	// 使用constant中的函数获取分类名称
	return constant.GetProductCategoryName(cateID)
}
