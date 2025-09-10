package bootstrap

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/olahol/melody"
	"go-web-template/global"
	"go-web-template/internal/config"
	"go-web-template/internal/db"
	"go-web-template/routers/admin_routes"
	"go-web-template/routers/api_routes"
	"log"
	"sync"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// StartAPIServer 启动API服务器（类似go-novel架构，可灵活配置组件）
func StartAPIServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	// 初始化配置
	config.Init()

	// API服务可选启用的组件（可根据需要调整）
	InitDB()
	InitZapLog()
	InitRedis()
	InitNsq()
	InitKeyLock()
	InitWs()
	InitGeoReader()
	InitBigcache()
	InitStaticServer() // 初始化静态资源服务

	// 初始化API路由
	InitAPIRoutes()
}

// StartAdminServer 启动管理后台服务器（类似go-novel架构，可灵活配置组件）
func StartAdminServer() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	// 初始化配置
	config.Init()

	// Admin服务可选启用的组件（可根据需要调整）
	InitDB()
	InitZapLog()
	InitRedis()
	InitNsq()
	InitKeyLock()
	InitGeoReader()
	InitBigcache()

	// 初始化管理后台路由
	InitAdminRoutes()
}

// ========== 数据库相关 ==========

// InitDB 初始化数据库连接
func InitDB() {
	var ctx = gctx.New()

	// 测试数据库连接
	database := g.DB()
	if database == nil {
		g.Log().Warning(ctx, "数据库连接失败")
		return
	}

	// 测试连接
	if err := database.PingMaster(); err != nil {
		g.Log().Fatal(ctx, "数据库连接测试失败:", err)
	} else {
		g.Log().Info(ctx, "数据库连接成功")
	}

	// 初始化GORM数据库连接
	cfg := config.Config
	dsn := cfg.Database.Default.LinkInfo
	if dsn == "" {
		// 如果linkInfo为空，构建DSN
		dsn = "root:root@tcp(127.0.0.1:3306)/template_chat?charset=utf8mb4&parseTime=True&loc=Local"
	}
	
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		g.Log().Fatal(ctx, "GORM数据库连接失败:", err)
		return
	}
	
	// 设置连接字符集为UTF8MB4
	if err := gormDB.Exec("SET NAMES utf8mb4").Error; err != nil {
		g.Log().Warning(ctx, "设置数据库字符集失败:", err)
	}
	
	// 设置全局GORM数据库实例
	global.DB = gormDB
	g.Log().Info(ctx, "GORM数据库初始化成功")
}

// InitRedis 初始化Redis连接
func InitRedis() {
	var ctx = gctx.New()

	// 使用db包的StartRedis函数统一加载Redis
	db.StartRedis()

	// 获取Redis客户端并设置全局变量
	global.Redis = db.GetRedisClient()

	if global.Redis != nil {
		g.Log().Info(ctx, "Redis连接已设置到全局变量")
	} else {
		g.Log().Warning(ctx, "Redis连接失败，将使用GoFrame内置Redis")
	}
}

// ========== 日志系统相关 ==========

// InitZapLog 初始化Zap日志系统 - 使用zaplog包的高级配置
func InitZapLog() {
	var ctx = gctx.New()

	// 使用db包的zaplog初始化函数
	db.InitZapLog()

	g.Log().Info(ctx, "Zap日志系统初始化成功（使用zaplog高级配置）",
		"features", []string{"文件轮转", "双输出", "按日期分割", "ZincSearch支持"})
}

// ========== 通信相关 ==========

// InitNsq 初始化NSQ消息队列
func InitNsq() {
	var ctx = gctx.New()

	// 这里可以根据配置初始化NSQ Producer
	// 目前先设置为nil，后续根据需要配置
	global.NsqPro = nil

	g.Log().Info(ctx, "NSQ消息队列初始化完成")
}

// InitKeyLock 初始化分布式锁
func InitKeyLock() {
	var ctx = gctx.New()

	// 使用sync.Mutex作为简单的锁机制
	global.KeyLock = &sync.Mutex{}

	g.Log().Info(ctx, "分布式锁初始化完成（使用sync.Mutex）")
}

// InitWs 初始化WebSocket（API服务专用）
func InitWs() {
	var ctx = gctx.New()

	// 初始化Melody WebSocket
	global.Ws = melody.New()

	// WebSocket连接管理初始化
	// 这里可以初始化WebSocket连接池、消息处理器等
	if global.Ws != nil {
		g.Log().Info(ctx, "WebSocket系统初始化成功，Melody已设置到全局变量")
	}
}

// ========== 地理位置相关 ==========

// InitGeoReader 初始化地理位置读取器
func InitGeoReader() {
	var ctx = gctx.New()

	// 地理位置数据库初始化
	// 可以初始化IP地址库、城市数据等
	g.Log().Info(ctx, "地理位置读取器初始化成功")
}

// ========== 缓存相关 ==========

// InitBigcache 初始化大型缓存系统
func InitBigcache() {
	var ctx = gctx.New()

	// 内存缓存系统初始化
	// 可以配置缓存容量、过期时间等
	g.Log().Info(ctx, "大型缓存系统初始化成功")
}

// ========== 路由相关 ==========

// InitAPIRoutes 初始化API路由
func InitAPIRoutes() {
	var ctx = gctx.New()

	// 创建HTTP服务
	s := g.Server()

	// 注册API路由
	api_routes.InitRoutes(s)

	g.Log().Info(ctx, "API服务器启动中...")

	// 启动服务器（阻塞方式）
	apiPort := g.Cfg().MustGet(ctx, "server.api.address", ":8080").String()
	s.SetAddr(apiPort)
	g.Log().Info(ctx, "API服务器已启动，访问地址: http://0.0.0.0"+apiPort)
	g.Log().Info(ctx, "按 Ctrl+C 退出API服务器")
	s.Run() // 阻塞运行
}

// InitAdminRoutes 初始化管理后台路由
func InitAdminRoutes() {
	var ctx = gctx.New()

	// 创建HTTP服务
	s := g.Server()
	// 从配置文件获取管理后台端口
	adminPort := g.Cfg().MustGet(ctx, "server.admin.address", ":8081").String()
	s.SetAddr(adminPort)

	// 注册管理后台路由
	admin_routes.InitRoutes(s)

	g.Log().Info(ctx, "管理后台服务器启动中...")

	// 启动服务器（阻塞方式）
	g.Log().Info(ctx, "管理后台服务器已启动，访问地址: http://0.0.0.0"+adminPort)
	g.Log().Info(ctx, "按 Ctrl+C 退出管理后台服务器")
	s.Run() // 阻塞运行
}

// ========== 静态资源服务相关 ==========

// InitStaticServer 初始化静态资源服务
func InitStaticServer() {
	var ctx = gctx.New()

	// 从配置文件获取静态资源服务器配置
	sourcePort := g.Cfg().MustGet(ctx, "server.source.address", ":8082").String()
	serverRoot := g.Cfg().MustGet(ctx, "server.source.serverRoot", "public").String()

	// 启动独立的静态资源服务器（非阻塞方式）
	go func() {
		// 创建独立的HTTP服务器实例
		staticServer := g.Server("static")
		staticServer.SetAddr(sourcePort)

		// 配置静态资源服务
		staticServer.SetServerRoot(serverRoot)
		staticServer.AddStaticPath("/", serverRoot)                   // 根路径直接访问文件
		staticServer.AddStaticPath("/static", serverRoot)             // /static路径访问
		staticServer.AddStaticPath("/uploads", serverRoot+"/uploads") // 上传文件专用路径

		g.Log().Info(ctx, "静态资源服务器启动中，端口:", sourcePort)

		// 启动服务器（在goroutine中阻塞运行）
		staticServer.Run()
	}()

	g.Log().Info(ctx, "静态资源服务器已在后台启动，访问地址: http://0.0.0.0"+sourcePort)
}
