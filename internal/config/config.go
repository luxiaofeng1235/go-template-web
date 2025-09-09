package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/joho/godotenv"
)

var Config *AppConfig

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Logs     LogsConfig     `yaml:"logs"`
	JWT      JWTConfig      `yaml:"jwt"`
	API      APIConfig      `yaml:"api"`
	OSS      OSSConfig      `yaml:"oss"`
}

type ServerConfig struct {
	API       APIServerConfig       `yaml:"api"`
	Admin     AdminServerConfig     `yaml:"admin"`
	Source    SourceServerConfig    `yaml:"source"`
	WebSocket WebSocketServerConfig `yaml:"websocket"`
}

type APIServerConfig struct {
	Address string `yaml:"address"`
	Name    string `yaml:"name"`
}

type AdminServerConfig struct {
	Address string `yaml:"address"`
	Name    string `yaml:"name"`
}

type SourceServerConfig struct {
	Address    string `yaml:"address"`
	ServerRoot string `yaml:"serverRoot"`
	Name       string `yaml:"name"`
}

type WebSocketServerConfig struct {
	Address string `yaml:"address"`
	Name    string `yaml:"name"`
}

type DatabaseConfig struct {
	Default DatabaseDefaultConfig `yaml:"default"`
}

type DatabaseDefaultConfig struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Pass            string `yaml:"pass"`
	Name            string `yaml:"name"`
	Charset         string `yaml:"charset"`
	Prefix          string `yaml:"prefix"`
	Debug           bool   `yaml:"debug"`
	LinkInfo        string `yaml:"linkInfo"`
	ConnMaxIdle     int    `yaml:"connMaxIdle"`
	ConnMaxOpen     int    `yaml:"connMaxOpen"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Address  string `yaml:"address"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

type LogsConfig struct {
	Level      int    `yaml:"level"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max-size"`
	MaxBackups int    `yaml:"max-backups"`
	MaxAge     int    `yaml:"max-age"`
	Compress   bool   `yaml:"compress"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

type APIConfig struct {
	Encrypt bool   `yaml:"encrypt"`
	AesKey  string `yaml:"aesKey"`
}

type OSSConfig struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	RoleArn         string `yaml:"roleArn"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"`
	ViewDomain      string `yaml:"viewDomain"`
	TokenExpireTime int    `yaml:"tokenExpireTime"`
	PolicyFile      string `yaml:"policyFile"`
}


// Init 初始化配置
func Init() {
	var ctx = gctx.New()
	
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		g.Log().Warning(ctx, "未找到.env文件，将使用默认配置:", err)
	} else {
		g.Log().Info(ctx, "成功加载.env环境变量")
	}
	
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

// GetOSS 获取OSS配置
func GetOSS() OSSConfig {
	return Config.OSS
}


// GetAPI 获取API配置
func GetAPI() APIConfig {
	return Config.API
}

// GetWebSocket 获取WebSocket配置
func GetWebSocket() WebSocketServerConfig {
	return Config.Server.WebSocket
}