# Go Web Template 开发规范

## 🚨 重要开发约束

**严格禁止在Controller中定义struct结构体！**
- 所有接口相关的请求/响应结构体必须在 `internal/models/` 对应的Model文件中定义
- Controller只负责参数解析、调用Service层和返回响应
- 保持代码层次清晰，结构体定义与业务逻辑分离

## 开发规范

### 0. 中文化规范

**项目要求全面中文化**：
- **注释**: 所有代码注释必须使用中文
- **变量命名**: 使用有意义的英文命名，但注释用中文说明
- **错误信息**: 所有用户可见的错误信息使用中文
- **日志输出**: 所有日志信息使用中文
- **文档**: 项目文档、README等全部使用中文
- **配置说明**: 配置文件中的注释使用中文

```go
// 示例：推荐的中文注释风格
type User struct {
    ID       int64     `gorm:"column:id" json:"id"`             // 用户ID
    Username string    `gorm:"column:username" json:"username"` // 用户名
    Email    string    `gorm:"column:email" json:"email"`       // 邮箱地址
    Status   int       `gorm:"column:status" json:"status"`     // 用户状态：1正常 0禁用
    CreateAt time.Time `gorm:"column:create_at" json:"create_at"` // 创建时间
}

// 用户登录
func (c *UserController) Login(r *ghttp.Request) {
    // 获取请求参数
    var req LoginReq
    if err := r.Parse(&req); err != nil {
        r.Response.WriteJson(g.Map{
            "code": constant.PARAM_ERROR,
            "msg":  "参数错误", // 中文错误信息
        })
        return
    }
    
    // 业务逻辑处理...
    global.Requestlog.Info("用户登录请求", "username", req.Username) // 中文日志
}
```

### 🏗️ 分层架构职责

**Controller层** - 接口控制器 (`api/controller/`)
- ✅ 负责：参数解析、参数验证、调用Service层、返回响应
- ❌ 禁止：定义结构体、编写业务逻辑、直接操作数据库
- 📁 位置：`api/controller/user.go`、`api/controller/product.go`

**Service层** - 业务逻辑层 (`internal/service/`)
- ✅ 负责：复杂业务逻辑、数据处理、事务管理、调用Model层
- ✅ 可以：密码加密、权限验证、数据转换、第三方API调用
- 📁 位置：`internal/service/user.go`、`internal/service/product.go`

**Model层** - 数据模型层 (`internal/models/`)
- ✅ 负责：结构体定义、简单数据查询、基础CRUD操作
- ❌ 禁止：复杂业务逻辑、跨表复杂查询、事务处理
- 📁 位置：`internal/models/UserModel.go`、`internal/models/ProductModel.go`

**Router层** - 路由配置层 (`routers/`)
- ✅ 负责：路由分组、中间件配置、接口路径定义
- 📁 位置：`routers/api_routes/`（前端接口）、`routers/admin_routes/`（后台管理）

### 1. 目录结构规范

- **模型文件**: `internal/models/` 下，命名格式：`XxxModel.go`
- **控制器文件**: `api/controller/` 下，按业务模块分文件
- **服务文件**: `internal/service/` 下，按业务模块分文件
- **路由文件**: `routers/` 下，分为 `api_routes` 和 `admin_routes`
- **工具函数**: `utils/` 下，按功能分文件

### 2. 模型定义规范

**严格禁止在Controller中定义结构体**：
- ❌ **禁止行为**: 在任何Controller文件中直接定义struct结构体
- ✅ **正确做法**: 所有结构体定义必须在 `internal/models/` 目录下的对应Model文件中

**模型文件职责**：
- 只定义结构体，不包含业务逻辑方法
- 包含数据库实体结构体
- 包含请求/响应结构体
- 包含接口相关的所有数据模型

```go
// 错误示例 ❌ - 禁止在Controller中定义
// api/controller/user.go
type LoginReq struct {  // 严格禁止这样做！
    Username string `json:"username"`
    Password string `json:"password"`
}

// 正确示例 ✅ - 在Model文件中定义
// internal/models/UserModel.go
package models

import "time"

// 数据库实体结构体
type User struct {
    ID       int64     `gorm:"column:id" json:"id"`             // 用户ID
    Username string    `gorm:"column:username" json:"username"` // 用户名
    Email    string    `gorm:"column:email" json:"email"`       // 邮箱地址
    Status   int       `gorm:"column:status" json:"status"`     // 用户状态：1正常 0禁用
    CreateAt time.Time `gorm:"column:create_at" json:"create_at"` // 创建时间
    UpdateAt time.Time `gorm:"column:update_at" json:"update_at"` // 更新时间
}

// 接口请求结构体 - 必须在Model中定义
type LoginReq struct {
    Username string `json:"username"` // 用户名
    Password string `json:"password"` // 密码
}

type RegisterReq struct {
    Username string `json:"username"` // 用户名
    Email    string `json:"email"`    // 邮箱
    Password string `json:"password"` // 密码
}

// 请求结构体使用组合模式
type UpdateUserReq struct {
    UserID int64 `json:"user_id"` // 用户ID
    RegisterReq                   // 组合注册请求结构体
}

// 响应结构体
type LoginRes struct {
    Token    string `json:"token"`     // JWT令牌
    UserInfo User   `json:"user_info"` // 用户信息
}
```

**Controller中的正确使用方式**：
```go
// api/controller/user.go
package controller

import (
    "go-web-template/internal/models"  // 导入models包
)

type UserController struct{}

// 用户登录 - 使用models中定义的结构体
func (c *UserController) Login(r *ghttp.Request) {
    // 使用models中定义的结构体
    var req models.LoginReq
    if err := r.Parse(&req); err != nil {
        // 处理错误...
    }
    
    // 业务逻辑处理...
    
    // 返回响应，使用models中定义的结构体
    res := models.LoginRes{
        Token:    "jwt_token",
        UserInfo: userInfo,
    }
    
    r.Response.WriteJson(res)
}
```

### 3. 字段定义规范

- **GORM字段定义**: 使用基本格式 `gorm:"column:name" json:"name"`
- **避免复杂约束**: 不在模型中定义复杂的数据库约束
- **不使用binding标签**: 移除所有 `binding:"required"` 标签

### 4. Service层设计规范

**基于go-novel的Service层开发标准**：

**函数设计模式**：
- ✅ 使用直接函数方式，避免过度面向对象设计
- ✅ 函数名清晰表达业务意图，如 `Info()`, `List()`, `Create()`, `Update()`
- ✅ 返回值统一使用 `(结果, error)` 模式

```go
// 推荐的Service函数设计
func Info(req *models.BookInfoReq) (rbook *models.BookInfoRes, err error) {
    // 业务逻辑实现
}

func List(req *models.BookListReq) (list []models.BookListRes, total int64, err error) {
    // 列表查询逻辑
}
```

**参数处理策略**：
- 📝 **复杂请求使用结构体**: 多个参数或复杂查询条件使用 `*models.XxxReq` 结构体
- 📝 **简单请求使用基本类型**: 单一参数查询直接使用 `int64`, `string` 等基本类型
- 📝 **立即参数验证**: 函数开头立即验证关键参数，使用中文错误信息

```go
// 复杂请求示例 - 使用结构体参数
func Info(req *models.BookInfoReq) (rbook *models.BookInfoRes, err error) {
    bookId := req.BookId
    if bookId <= 0 {
        err = fmt.Errorf("%v", "小说ID为空")
        return
    }
    
    userId := req.UserId
    if userId <= 0 {
        err = fmt.Errorf("%v", "用户ID为空")
        return
    }
    
    // 业务逻辑...
}

// 简单请求示例 - 直接使用基本类型
func GetUserById(userId int64) (user *models.User, err error) {
    if userId <= 0 {
        err = fmt.Errorf("%v", "用户ID无效")
        return
    }
    
    // 数据查询...
}
```

**错误处理规范**：
- 🈯 **中文错误信息**: 所有用户可见错误必须使用中文
- 🈯 **立即返回**: 参数验证失败立即返回，避免继续执行
- 🈯 **错误信息具体化**: 错误信息要具体说明问题，便于调试

**日志记录模式**：
- 📊 **使用分类日志器**: 根据业务模块选择合适的日志器 (`global.Sqllog`, `global.Requestlog` 等)
- 📊 **结构化日志**: 使用键值对方式记录关键信息
- 📊 **中文日志信息**: 所有日志输出使用中文，便于运维人员理解

```go
// 日志记录示例
func Info(req *models.BookInfoReq) (rbook *models.BookInfoRes, err error) {
    // 记录请求日志
    global.Requestlog.Info("查询小说信息", "bookId", req.BookId, "userId", req.UserId)
    
    // 数据库操作前记录
    global.Sqllog.Info("开始查询小说基本信息", "bookId", req.BookId)
    
    // 业务逻辑...
    
    // 结果日志
    global.Requestlog.Info("小说信息查询成功", "bookId", req.BookId, "title", book.BookName)
    
    return
}
```

**数据库操作模式**：
- 🗄️ **GORM链式操作**: 使用GORM的链式查询，便于构建复杂查询条件
- 🗄️ **开启Debug模式**: 开发环境开启 `.Debug()` 模式，便于调试SQL
- 🗄️ **分步骤查询**: 复杂业务逻辑分多个查询步骤，保证代码可读性

```go
// GORM查询示例
func Info(req *models.BookInfoReq) (rbook *models.BookInfoRes, err error) {
    bookId := req.BookId
    if bookId <= 0 {
        err = fmt.Errorf("%v", "小说ID为空")
        return
    }
    
    var book models.Book
    // 使用GORM链式查询 + Debug模式
    err = global.DB.Debug().Where("id = ?", bookId).First(&book).Error
    if err != nil {
        global.Errlog.Error("查询小说失败", "bookId", bookId, "error", err)
        return
    }
    
    // 数据转换
    rbook = &models.BookInfoRes{
        BookId:   book.ID,
        BookName: book.BookName,
        Author:   book.Author,
        // 其他字段...
    }
    
    return
}
```

### 5. 常量和枚举管理规范

**常量定义位置**: `internal/constant/` 目录，按业务模块分文件
- 避免硬编码，所有常量都在constant包中统一定义
- 按业务领域分文件：`product.go`、`file.go`、`user.go`、`chat_common.go`等
- Service层和组件中经常使用这些常量

**枚举定义模式**: 使用结构体+切片方式，便于前端遍历
```go
// internal/constant/product.go
package constant

// 产品分类结构体
type ProductCategory struct {
    Value int    `json:"value"` // 分类ID  
    Label string `json:"label"` // 分类名称
}

// 产品分类枚举列表（前端下拉框直接遍历）
var ProductCategoryList = []ProductCategory{
    {Value: 0, Label: "请选择"},
    {Value: 1, Label: "AI助手"}, 
    {Value: 2, Label: "机器学习"},
    {Value: 3, Label: "计算机视觉"},
}

// 产品状态常量
const (
    PRODUCT_STATUS_NORMAL  = 1 // 正常
    PRODUCT_STATUS_DISABLE = 0 // 禁用
)
```

**Service层使用常量示例**:
```go
// internal/service/product.go
import "go-web-template/internal/constant"

func (s *productService) CreateProduct(ctx context.Context, req *models.CreateProductReq) error {
    // 使用常量避免硬编码
    if req.CategoryID <= 0 || req.CategoryID >= len(constant.ProductCategoryList) {
        return errors.New("产品分类无效")
    }
    
    product := &models.Product{
        Name:       req.Name,
        CategoryID: req.CategoryID,
        Status:     constant.PRODUCT_STATUS_NORMAL, // 使用常量
    }
    
    return global.DB.Create(product).Error
}
```

**常量文件组织**:
- `all_const.go` - 通用常量（HTTP状态码、响应消息等）
- `product.go` - 产品相关常量和枚举
- `user.go` - 用户相关常量和枚举  
- `file.go` - 文件类型相关常量
- `chat_common.go` - 聊天相关常量

## 标准开发流程

### 🔄 开发流程步骤

**第一步：定义结构体** (`internal/models/XxxModel.go`)
```go
// internal/models/UserModel.go
type LoginReq struct {
    Username string `json:"username"` // 用户名
    Password string `json:"password"` // 密码
}

type LoginRes struct {
    Token    string `json:"token"`     // JWT令牌
    UserInfo User   `json:"user_info"` // 用户信息
}
```

**第二步：编写接口函数** (`api/controller/` 或对应业务controller)
```go
// api/controller/user.go
func (c *UserController) Login(r *ghttp.Request) {
    // 参数解析
    var req models.LoginReq
    if err := r.Parse(&req); err != nil {
        utils.ParamError(r, "参数错误")
        return
    }
    
    // 调用Service层处理业务逻辑
    result, err := service.User.Login(r.Context(), req.Username, req.Password)
    if err != nil {
        utils.Fail(r, err, "登录失败")
        return
    }
    
    utils.Success(r, result, "登录成功")
}
```

**第三步：配置分组路由** (`routers/api_routes/` 或 `routers/admin_routes/`)
```go
// routers/api_routes/user.go
func InitUserRoutes(group *ghttp.RouterGroup) {
    userCtrl := &controller.UserController{}
    
    group.POST("/login", userCtrl.Login)        // 用户登录
    group.POST("/register", userCtrl.Register)  // 用户注册
    group.GET("/profile", userCtrl.GetProfile)  // 获取用户信息
}
```

**第四步：编写Service业务逻辑** (`internal/service/`)
```go
// internal/service/user.go
func (s *userService) Login(ctx context.Context, username, password string) (*models.LoginRes, error) {
    // 业务逻辑处理
    user, err := s.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, err
    }
    
    // 密码验证等复杂逻辑
    if !s.VerifyPassword(password, user.Password) {
        return nil, errors.New("用户名或密码错误")
    }
    
    // 生成JWT等
    token := s.GenerateToken(user.ID)
    
    return &models.LoginRes{
        Token:    token,
        UserInfo: *user,
    }, nil
}
```

**第五步：Model层提供数据查询** (`internal/models/`)
```go
// Model层只提供简单的数据查询，不包含复杂业务逻辑
type User struct {
    ID       int64  `gorm:"column:id" json:"id"`
    Username string `gorm:"column:username" json:"username"`
    Password string `gorm:"column:password" json:"-"` // 密码不返回给前端
    Email    string `gorm:"column:email" json:"email"`
}

// 简单的查询方法（可以在Model中定义基础查询）
func GetUserByUsername(username string) (*User, error) {
    var user User
    err := global.DB.Where("username = ?", username).First(&user).Error
    return &user, err
}
```

### 📁 路由分组规范

**API路由** - 用于前端接口 (`routers/api_routes/`)
```go
// routers/api_routes/routes.go
func InitRoutes(s *ghttp.Server) {
    // API分组
    apiGroup := s.Group("/api")
    
    // 用户相关路由
    InitUserRoutes(apiGroup.Group("/user"))
    
    // 产品相关路由  
    InitProductRoutes(apiGroup.Group("/product"))
}
```

**Admin路由** - 用于后台管理 (`routers/admin_routes/`)
```go
// routers/admin_routes/routes.go
func InitRoutes(s *ghttp.Server) {
    // Admin分组
    adminGroup := s.Group("/admin")
    
    // 用户管理路由
    InitUserRoutes(adminGroup.Group("/user"))
    
    // 产品管理路由
    InitProductRoutes(adminGroup.Group("/product"))
}
```

## WebSocket开发规范

### 推荐使用Gorilla WebSocket

项目推荐使用 `github.com/gorilla/websocket` 作为WebSocket实现：

**安装依赖**:
```bash
go get github.com/gorilla/websocket
```

**基础使用示例**:
```go
// WebSocket升级器配置
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许跨域，生产环境需要严格控制
    },
}

// WebSocket处理函数
func handleWebSocket(r *ghttp.Request) {
    // 升级HTTP连接为WebSocket
    conn, err := upgrader.Upgrade(r.Response.ResponseWriter, r.Request, nil)
    if err != nil {
        global.Wslog.Error("WebSocket升级失败:", err)
        return
    }
    defer conn.Close()

    // 消息处理循环
    for {
        // 读取消息
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            global.Wslog.Error("读取消息失败:", err)
            break
        }

        // 处理消息逻辑
        global.Wslog.Info("收到消息:", string(message))

        // 回复消息
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            global.Wslog.Error("发送消息失败:", err)
            break
        }
    }
}
```

**路由注册**:
```go
// 在路由文件中注册WebSocket接口
func InitWebSocketRoutes(group *ghttp.RouterGroup) {
    group.GET("/ws", handleWebSocket) // WebSocket连接端点
}
```

**推荐原因**:
- Go生态系统标准WebSocket库
- 广泛的生产环境验证
- 完整的WebSocket协议支持
- 与GoFrame完美集成
- 丰富的文档和示例

## 高级日志系统

### Zap日志特性

基于go-novel的zaplog实现，提供以下特性：

- **日志轮转**: 支持按大小、时间自动轮转
- **双输出**: 同时输出到文件和控制台
- **按日期分割**: 日志文件按日期自动命名 (20250909.log)
- **模块化分类**: 14个专用日志记录器
- **ZincSearch支持**: 可选的日志搜索和分析功能

### 日志记录器分类

- `global.Errlog` - 系统错误日志
- `global.Sqllog` - 数据库SQL执行日志
- `global.Requestlog` - HTTP请求日志
- `global.Paylog` - 支付相关日志
- `global.Wslog` - WebSocket连接日志
- `global.Nsqlog` - 消息队列日志
- `global.Collectlog` - 数据采集日志
- `global.Updatelog` - 数据更新日志
- 其他业务模块专用日志记录器

### 日志使用示例

```go
// 错误日志
global.Errlog.Error("用户登录失败", "username", username, "error", err)

// SQL日志
global.Sqllog.Info("执行SQL查询", "sql", sqlStr, "duration", time.Since(start))

// 请求日志
global.Requestlog.Info("API请求", "method", r.Method, "url", r.URL.Path, "ip", r.GetClientIp())

// WebSocket日志
global.Wslog.Info("WebSocket连接建立", "clientId", clientId, "remoteAddr", conn.RemoteAddr())
```

## 启动方式

```bash
# 启动API服务器
go run api.go

# 启动管理后台服务器
go run admin.go
```

## 注意事项

### 📋 开发流程约束

1. **结构体定义**: **严格禁止在Controller中定义struct**，所有结构体必须在`internal/models/`中定义
2. **分层职责**: 严格按照 Controller → Service → Model 的调用顺序，不允许跨层调用
3. **Model层限制**: Model层不写复杂业务逻辑，只提供基础的数据查询操作
4. **Service层职责**: 所有复杂业务逻辑都在Service层实现，通过Service调用Model查询
5. **路由分组**: 根据业务场景选择API路由(前端接口)或Admin路由(后台管理)
6. **接口定义**: 先在Controller中定义接口函数，再配置到对应的路由分组中

### 🛠️ 技术约束

7. **常量管理**: 严禁硬编码，所有常量都在`internal/constant/`中按业务模块定义
8. **配置管理**: 所有配置通过yaml统一管理，避免硬编码
9. **资源分离**: API/Admin专注业务逻辑，静态资源统一通过8082端口
10. **日志管理**: 使用zap结构化日志，支持日志轮转和压缩
11. **全局变量**: 统一通过global包管理，避免循环引用
12. **启动封装**: 参考go-novel模式，通过db包封装启动逻辑

### 📂 常量使用规范

- **禁止硬编码**: 代码中不允许出现魔法数字或字符串常量
- **按业务分组**: 产品相关常量放在`product.go`，用户相关放在`user.go`
- **Service层优先使用**: Service层和各种组件优先使用constant包中的常量
- **枚举结构化**: 使用结构体+切片模式定义枚举，便于前端遍历使用

### 🚫 严格禁止的行为

```go
// ❌ 绝对禁止在Controller中定义结构体
// api/controller/user.go
type LoginRequest struct {     // 这样做是违规的！
    Username string `json:"username"`
}

// ❌ 绝对禁止在Controller中定义任何struct
type UserResponse struct {     // 这样做也是违规的！
    Message string `json:"message"`
}
```

### ✅ 正确的做法

```go
// ✅ 在Model文件中定义所有结构体
// internal/models/UserModel.go
type LoginRequest struct {
    Username string `json:"username"` // 用户名
    Password string `json:"password"` // 密码
}

// ✅ 在Controller中导入并使用
// api/controller/user.go
import "go-web-template/internal/models"

func (c *UserController) Login(r *ghttp.Request) {
    var req models.LoginRequest  // 正确使用方式
    // ...
}
```

## 版本信息

- **创建日期**: 2025-09-09
- **最后更新**: 2025-09-09