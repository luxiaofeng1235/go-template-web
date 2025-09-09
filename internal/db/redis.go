package db

import (
	"context"
	"fmt"
	"log"

	"go-web-template/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
)

var RedisClient *redis.Client

// GetRedisConfig 获取Redis配置（参考go-novel的GetRedis，但不区分环境变量）
func GetRedisConfig() (addr string, passwd string, defaultdb int) {
	// 从配置文件读取Redis配置
	cfg := config.Config
	if cfg != nil {
		addr = cfg.Redis.Default.Address
		passwd = cfg.Redis.Default.Password
		defaultdb = cfg.Redis.Default.DB
	} else {
		// 默认配置（配置文件加载失败时的备用）
		addr = "127.0.0.1:6379"
		passwd = ""
		defaultdb = 0
	}
	return
}

// InitRedis 初始化Redis（参考go-novel的InitRedis）
func InitRedis(addr string, passwd string, defaultdb int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       defaultdb,
	})

	var ctx = context.Background()
	err := redisClient.Ping(ctx).Err()
	if err != nil || redisClient == nil {
		log.Println(fmt.Sprintf("初始化Redis异常：%v，将使用GoFrame内置Redis", err))
		return nil
	} else {
		log.Println("Redis连接成功")
	}

	RedisClient = redisClient
	return redisClient
}

// StartRedis 启动Redis连接（参考go-novel但使用GoFrame配置）
func StartRedis() {
	addr, passwd, defaultdb := GetRedisConfig()

	// 先尝试使用go-redis客户端
	client := InitRedis(addr, passwd, defaultdb)

	// 如果go-redis连接失败，GoFrame的g.Redis()会自动使用配置文件连接
	if client == nil {
		var ctx = context.Background()
		_, err := g.Redis().Do(ctx, "PING")
		if err != nil {
			log.Printf("Redis连接失败: %v", err)
		} else {
			log.Println("使用GoFrame Redis连接成功")
		}
	}
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	return RedisClient
}
