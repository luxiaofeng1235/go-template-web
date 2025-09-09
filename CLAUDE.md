# Go Web Template 项目介绍

## 项目概述

本项目是基于 GoFrame v2 框架的私域聊天平台，参考 go-novel 的架构设计，采用分层设计模式，支持 API 服务器、管理后台和静态资源服务的分离部署。

### 主要功能需求

**访问控制**：
- 进入网页必须输入正确的密钥/密码/暗号
- 密钥/密码/暗号可配置且可设置有效期
- 实现私域个人站点访问控制

**核心功能**：
- 🤖 **AI搜索** - 智能搜索功能
- 🎨 **AI生图** - AI图像生成功能  
- 📹 **AI视频** - AI视频生成功能
- 💬 **群聊** - 广播群聊，无需建群，使用相同密钥的用户自动组成群聊
- 🔒 **私聊** - 用户间私人聊天功能
- 🎙️ **音视频聊天** - 支持实时音频和视频通话，通过WebSocket实现信令控制

**设备兼容**：
- PC和移动端兼容
- 纯网页形式，无需安装APP

### 用户工作流程

1. **用户访问** → 打开网站，前端JS自动收集设备信息
2. **设备指纹** → 前端计算设备指纹并发送给后端  
3. **身份生成** → 后端基于设备指纹生成固定的 access_key
4. **建立连接** → 用户获得固定身份标识，可以开始聊天交互

通过设备指纹技术，确保每个设备都有唯一且固定的身份标识，实现免注册的稳定用户体系。

**## 技术栈

- **框架**: GoFrame v2
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **日志**: Zap Logger (高级日志轮转)
- **消息队列**: NSQ
- **WebSocket**: Gorilla WebSocket (推荐标准库)
- **音视频通信**: WebRTC + WebSocket信令
- **配置管理**: YAML**


## 架构设计

### 服务器架构

项目采用三服务器分离架构：

1. **API服务器** (`:8080`) - 处理前端API请求
2. **管理后台服务器** (`:8081`) - 处理后台管理功能
3. **静态资源服务器** (`:8082`) - 提供图片等静态资源，在bootstrap中后台启动

### 启动流程

```
启动文件 -> db包 -> server包 -> bootstrap包 -> 具体初始化
```

**调用链路:**
- `api.go` -> `db.StartAPIServer()` -> `server.StartAPIServer()` -> `bootstrap.StartAPIServer()`
- `admin.go` -> `db.StartAdminServer()` -> `server.StartAdminServer()` -> `bootstrap.StartAdminServer()`

### 全局变量管理

参考go-novel的global.go设计，统一管理全局资源：

- `global.DB` - GORM数据库连接
- `global.Redis` - Redis客户端
- `global.NsqPro` - NSQ消息队列
- `global.KeyLock` - 分布式锁
- `global.Ws` - WebSocket连接 (使用Gorilla WebSocket)
- 各种专用日志记录器：`Errlog`、`Sqllog`、`Requestlog`等

## 核心特性

### 🎯 分层架构
- **Controller层**: 接口控制，参数处理
- **Service层**: 业务逻辑，数据处理
- **Model层**: 数据模型，结构体定义
- **Router层**: 路由配置，分组管理

### 📊 高级日志系统
- **基于go-novel的zaplog**: 从go-novel迁移的成熟日志解决方案
- **日志轮转**: 支持按大小、时间自动轮转
- **双输出**: 同时输出到文件和控制台
- **按日期分割**: 日志文件按日期自动命名
- **模块化分类**: 14个专用日志记录器，按业务模块分离
- **ZincSearch支持**: 可选的日志搜索和分析功能

### 🌐 WebSocket支持
- **Gorilla WebSocket**: 使用Go生态标准WebSocket库
- **完整协议支持**: 支持所有WebSocket特性
- **高性能**: 经过大量生产环境验证
- **易于集成**: 与GoFrame完美集成

### ⚙️ 配置管理
- **YAML配置**: 统一的配置文件管理
- **多环境支持**: 支持开发、测试、生产环境配置
- **热重载**: 支持配置文件热重载

## 配置文件结构

```yaml
# 服务器配置
server:
  api:
    address: ":8080"
    name: "API服务器"
  admin:
    address: ":8081"
    name: "后台管理服务器"
  source:
    address: ":8082"
    serverRoot: "public"
    name: "静态资源服务器"

# 数据库配置
database:
  default:
    type: mysql
    host: 127.0.0.1
    port: 3306
    user: root
    pass: root
    name: template_chat
    charset: utf8mb4
    prefix: ls_
    debug: true

# Redis配置
redis:
  default:
    address: "127.0.0.1:6379"
    db: 0
    password: ""

# 日志配置
logs:
  level: -1  # -1:Debug, 0:Info, 1:Warn, 2:Error
  path: logs
  max-size: 50      # 文件最大大小(MB)
  max-backups: 100  # 备份数
  max-age: 30       # 存放时间(天)
  compress: false   # 是否压缩
```

## 快速开始

### 环境要求
- Go 1.23.0+
- MySQL 5.7+
- Redis 6.0+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd go-web-template
```

2. **安装依赖**
```bash
go mod tidy
```

3. **连接数据库**
```bash
# 数据库已存在，直接连接查看
USE template_chat;
SHOW TABLES;  # 查看现有数据表
```

4. **创建必要目录**
```bash
mkdir -p logs public/uploads
```

5. **启动服务**
```bash
# 启动API服务器
go run api.go

# 启动管理后台服务器  
go run admin.go
```

### 访问地址
- **API服务器**: http://localhost:8080
- **管理后台**: http://localhost:8081
- **静态资源**: http://localhost:8082

## 数据库设计

项目包含以下主要数据表：

- `ls_user` - 用户表
- `ls_product` - 产品表
- `ls_ai_chat_log` - AI聊天记录
- `ls_ai_work` - AI工作任务
- `ls_chat_message` - 聊天消息
- `ls_config` - 系统配置
- `ls_file_cate` - 文件分类
- `ls_file` - 文件表
- `ls_meeting_room` - 会议室
- `ls_secret_key` - 密钥管理
- `ls_user_note` - 用户笔记
- `ls_user_session` - 用户会话

## 开发规范

详细的开发规范请查看 [CLAUDE_DEV.md](./CLAUDE_DEV.md) 文档。

## 版本信息

- **Go版本**: 1.23.0
- **GoFrame版本**: v2.9.3
- **WebSocket库**: github.com/gorilla/websocket (推荐)
- **日志系统**: 基于go-novel的zaplog高级实现
- **创建日期**: 2025-09-09
- **最后更新**: 2025-09-09