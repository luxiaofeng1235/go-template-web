package service

import (
	"context"
	"errors"
	"go-web-template/internal/constant"
	"go-web-template/internal/models"

	"github.com/gogf/gf/v2/frame/g"
)

type productService struct{}

var Product = productService{}

// GetList 获取产品列表
func (s productService) GetList(ctx context.Context, page, pageSize int) ([]models.ProductModel, int64, error) {
	var products []models.ProductModel

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询产品列表
	err := g.Model(constant.TABLE_PRODUCT).
		Where("status", constant.PRODUCT_STATUS_NORMAL).
		OrderDesc("created_at").
		Limit(offset, pageSize).
		Scan(&products)
	if err != nil {
		return nil, 0, err
	}

	// 查询总数
	total, err := g.Model(constant.TABLE_PRODUCT).
		Where("status", constant.PRODUCT_STATUS_NORMAL).
		Count()
	if err != nil {
		return nil, 0, err
	}

	return products, int64(total), nil
}

// GetDetail 获取产品详情
func (s productService) GetDetail(ctx context.Context, id int) (*models.ProductModel, error) {
	var product models.ProductModel
	err := g.Model(constant.TABLE_PRODUCT).
		Where("id", id).
		Where("status", constant.PRODUCT_STATUS_NORMAL).
		Scan(&product)
	if err != nil {
		return nil, err
	}
	if product.ID == 0 {
		return nil, errors.New("产品不存在")
	}

	return &product, nil
}

// Create 创建产品
func (s productService) Create(ctx context.Context, name, description string, price float64) error {
	_, err := g.Model(constant.TABLE_PRODUCT).Data(g.Map{
		"name":        name,
		"description": description,
		"price":       price,
		"status":      constant.PRODUCT_STATUS_NORMAL,
	}).Insert()

	return err
}

// Update 更新产品
func (s productService) Update(ctx context.Context, id int, data g.Map) error {
	result, err := g.Model(constant.TABLE_PRODUCT).
		Where("id", id).
		Where("status", constant.PRODUCT_STATUS_NORMAL).
		Data(data).
		Update()
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("产品不存在或更新失败")
	}

	return nil
}

// Delete 删除产品（软删除）
func (s productService) Delete(ctx context.Context, id int) error {
	result, err := g.Model(constant.TABLE_PRODUCT).
		Where("id", id).
		Where("status", constant.PRODUCT_STATUS_NORMAL).
		Data("status", constant.PRODUCT_STATUS_DISABLE).
		Update()
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("产品不存在或删除失败")
	}

	return nil
}
