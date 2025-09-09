package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var Config *AppConfig

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Logger   LoggerConfig   `yaml:"logger"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Address    string `yaml:"address"`
	ServerRoot string `yaml:"serverRoot"`
}

type DatabaseConfig struct {
	Type     string            `yaml:"type"`
	Hostname string            `yaml:"hostname"`
	Database string            `yaml:"database"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
	Hostport int               `yaml:"hostport"`
	Debug    bool              `yaml:"debug"`
	Charset  string            `yaml:"charset"`
	Prefix   string            `yaml:"prefix"`
	Default  map[string]interface{} `yaml:"default"`
}

type RedisConfig struct {
	Default map[string]interface{} `yaml:"default"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

// Init 初始化配置
func Init() {
	var ctx = gctx.New()
	Config = &AppConfig{}
	
	// 从配置文件加载配置
	if err := g.Cfg().MustGet(ctx, "").Struct(Config); err != nil {
		g.Log().Fatal(ctx, "配置文件加载失败:", err)
	}
	
	g.Log().Info(ctx, "配置加载成功")
}

// GetDB 获取数据库配置
func GetDB() DatabaseConfig {
	return Config.Database
}

// GetJWT 获取JWT配置
func GetJWT() JWTConfig {
	return Config.JWT
}