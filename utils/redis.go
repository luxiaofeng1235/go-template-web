package utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go-web-template/internal/constant"
)

// Redis 单独的Redis工具函数

// GetRedis 获取Redis实例
func GetRedis() *gredis.Redis {
	return g.Redis()
}

// SetToken 设置Token缓存
func SetToken(ctx context.Context, token string, userID int, expireSeconds ...int) error {
	expire := 86400 // 默认24小时
	if len(expireSeconds) > 0 {
		expire = expireSeconds[0]
	}
	
	key := constant.REDIS_TOKEN_PREFIX + token
	err := g.Redis().SetEX(ctx, key, userID, int64(expire))
	return err
}

// GetTokenUserID 根据Token获取用户ID
func GetTokenUserID(ctx context.Context, token string) (int, error) {
	key := constant.REDIS_TOKEN_PREFIX + token
	result, err := g.Redis().Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return result.Int(), nil
}

// DeleteToken 删除Token
func DeleteToken(ctx context.Context, token string) error {
	key := constant.REDIS_TOKEN_PREFIX + token
	_, err := g.Redis().Del(ctx, key)
	return err
}

// SetUserCache 设置用户缓存
func SetUserCache(ctx context.Context, userID int, userData interface{}, expireSeconds ...int) error {
	expire := 3600 // 默认1小时
	if len(expireSeconds) > 0 {
		expire = expireSeconds[0]
	}
	
	key := constant.REDIS_USER_PREFIX + fmt.Sprintf("%d", userID)
	err := g.Redis().SetEX(ctx, key, userData, int64(expire))
	return err
}

// GetUserCache 获取用户缓存
func GetUserCache(ctx context.Context, userID int) (interface{}, error) {
	key := constant.REDIS_USER_PREFIX + fmt.Sprintf("%d", userID)
	return g.Redis().Get(ctx, key)
}

// SetProductCache 设置产品缓存
func SetProductCache(ctx context.Context, productID int, productData interface{}, expireSeconds ...int) error {
	expire := 1800 // 默认30分钟
	if len(expireSeconds) > 0 {
		expire = expireSeconds[0]
	}
	
	key := constant.REDIS_PRODUCT_PREFIX + fmt.Sprintf("%d", productID)
	err := g.Redis().SetEX(ctx, key, productData, int64(expire))
	return err
}

// GetProductCache 获取产品缓存
func GetProductCache(ctx context.Context, productID int) (interface{}, error) {
	key := constant.REDIS_PRODUCT_PREFIX + fmt.Sprintf("%d", productID)
	return g.Redis().Get(ctx, key)
}

// TestRedisConnection 测试Redis连接（单独启动时使用）
func TestRedisConnection() error {
	ctx := gctx.New()
	_, err := g.Redis().Do(ctx, "PING")
	return err
}

// FlushRedis 清空Redis缓存（开发调试用）
func FlushRedis(ctx context.Context) error {
	err := g.Redis().FlushAll(ctx)
	return err
}