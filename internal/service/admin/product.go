/*
 * @file: product.go
 * @description: 商品管理相关业务逻辑处理 - Admin后台管理层服务
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package admin

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"
	"go-web-template/utils"
	"strings"
)

// GetProductList 获取商品列表 - 后台管理专用
// @param req *models.ProductListReq 商品列表查询请求参数
// @return list []models.ProductListItem 商品列表数据
// @return total int64 总记录数
// @return err error 错误信息
func GetProductList(req *models.ProductListReq) (list []models.ProductListItem, total int64, err error) {
	// 参数默认值处理
	if req.PageNo <= 0 {
		req.PageNo = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	query := global.DB.Model(&models.Product{})

	// 添加分类筛选条件
	if req.CateID > 0 {
		query = query.Where("cate_id = ?", req.CateID)
	}

	// 获取总记录数
	err = query.Count(&total).Error
	if err != nil {
		global.Sqllog.Errorf("查询商品总数失败 err=%v", err.Error())
		return
	}

	// 分页查询
	offset := (req.PageNo - 1) * req.PageSize
	var products []models.Product

	err = query.Select("id, product_name, cate_id, intro, logo, qrcode, created_at, updated_at").
		Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&products).Error

	if err != nil {
		global.Sqllog.Errorf("查询后台商品列表失败 err=%v", err.Error())
		return
	}

	// 转换为返回格式
	list = make([]models.ProductListItem, 0, len(products))

	for _, product := range products {
		item := models.ProductListItem{
			ID:          product.ID,
			ProductName: product.ProductName,
			CateID:      product.CateID,
			Intro:       product.Intro,
			Logo:        processAdminLogoURL(product.Logo),
			QRCode:      product.QRCode,
			CateName:    getCategoryName(product.CateID),
		}
		list = append(list, item)
	}

	return list, total, nil
}

// GetCategoryList 获取商品分类列表 - 后台管理专用
// @return categories []constant.ProductCategory 分类列表数据
// @return err error 错误信息
func GetCategoryList() (categories []constant.ProductCategory, err error) {
	return constant.ProductCategoryList, nil
}

// processAdminLogoURL 处理后台管理的Logo URL
// @param logo string 原始Logo路径
// @return string 处理后的Logo URL
func processAdminLogoURL(logo string) string {
	if logo == "" {
		return ""
	}

	// 如果已经是完整URL，直接返回
	if len(logo) > 7 && (logo[:7] == "http://" || logo[:8] == "https://") {
		return logo
	}

	// 后台管理统一使用静态资源服务器地址
	baseURL := "http://localhost:8082"

	// 如果是相对路径，拼接域名
	if logo[0] == '/' {
		return baseURL + logo
	}

	// 其他情况，添加/前缀后拼接
	return baseURL + "/" + logo
}

// getCategoryName 根据分类ID获取分类名称
// @param cateID int 分类ID
// @return string 分类名称
func getCategoryName(cateID int) string {
	return constant.GetProductCategoryName(cateID)
}

// SaveProduct 保存商品信息 - 支持新增和编辑
// @param req *models.SaveProductReq 保存商品请求参数
// @return result *models.SaveProductRes 保存结果数据
// @return err error 错误信息
func SaveProduct(req *models.SaveProductReq) (result *models.SaveProductRes, err error) {
	// 字段预处理 - 使用strings.TrimSpace处理字符串字段
	productName := strings.TrimSpace(req.ProductName)
	intro := strings.TrimSpace(req.Intro)
	logo := strings.TrimSpace(req.Logo)
	qrcode := strings.TrimSpace(req.QRCode)

	// 基本参数验证
	if productName == "" {
		err = fmt.Errorf("产品名称必须输入")
		return
	}

	// 检查商品名称是否重复
	query := global.DB.Where("product_name = ?", productName)
	if req.ID > 0 {
		query = query.Where("id != ?", req.ID)
	}

	var existingProduct models.Product
	err = query.First(&existingProduct).Error
	if err == nil {
		// 找到了重复的商品名称
		err = fmt.Errorf("产品名称已存在，请更换产品名称")
		return
	}

	// 如果错误不是"记录未找到"，说明是数据库查询错误
	if err.Error() != "record not found" {
		return
	}

	// 清除错误，继续执行（没找到重复记录是正常的）
	err = nil

	now := utils.GetUnix()

	// 判断是新增还是编辑
	if req.ID > 0 {
		// 编辑操作
		updates := map[string]interface{}{
			"product_name": productName,
			"cate_id":      req.CateID,
			"intro":        intro,
			"logo":         logo,
			"qrcode":       qrcode,
			"updated_at":   now,
		}

		err = global.DB.Model(&models.Product{}).Where("id = ?", req.ID).Updates(updates).Error
		if err != nil {
			global.Sqllog.Errorf("编辑商品失败 err=%v", err.Error())
			return
		}

		result = &models.SaveProductRes{
			ID:          req.ID,
			ProductName: productName,
			CateID:      req.CateID,
			CreatedAt:   0, // 编辑时不返回创建时间
			UpdatedAt:   now,
		}
	} else {
		// 新增操作
		product := models.Product{
			ProductName: productName,
			CateID:      req.CateID,
			Intro:       intro,
			Logo:        logo,
			QRCode:      qrcode,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		err = global.DB.Create(&product).Error
		if err != nil {
			global.Sqllog.Errorf("新增商品失败 err=%v", err.Error())
			return
		}

		result = &models.SaveProductRes{
			ID:          product.ID,
			ProductName: product.ProductName,
			CateID:      product.CateID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	return result, nil
}

// DeleteProduct 删除商品
// @param productID int64 商品ID
// @return err error 错误信息
func DeleteProduct(productID int64) (err error) {
	if productID <= 0 {
		return fmt.Errorf("商品ID不能为空")
	}

	// 执行删除操作
	err = global.DB.Where("id = ?", productID).Delete(&models.Product{}).Error
	if err != nil {
		global.Sqllog.Errorf("删除商品失败 err=%v", err.Error())
		return
	}

	return nil
}
