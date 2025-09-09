/*
 * @file: global.go
 * @description: 全局变量定义和管理 - 参考go-novel的global设计模式
 * @author: go-web-template
 * @created: 2025-09-09
 * @version: 1.0.0
 * @license: MIT License
 */

package global

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/nsqio/go-nsq"
	"github.com/olahol/melody"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
全局变量定义 - 统一管理项目中的共享资源
参考go-novel的架构设计，将所有全局资源集中管理，避免循环引用
*/
var (
	// === 数据层相关 ===
	DB    *gorm.DB      // GORM数据库连接实例
	Redis *redis.Client // Redis缓存客户端

	// === 消息和通信相关 ===
	NsqPro  *nsq.Producer  // NSQ消息队列生产者
	KeyLock *sync.Mutex    // 分布式锁（使用sync.Mutex作为简单实现）
	Ws      *melody.Melody // WebSocket连接管理器

	// === 日志系统相关 ===
	// 基础日志记录器
	Errlog     *zap.SugaredLogger // 系统错误日志
	Sqllog     *zap.SugaredLogger // 数据库SQL执行日志
	Requestlog *zap.SugaredLogger // HTTP请求日志
	Paylog     *zap.SugaredLogger // 支付相关日志

	// WebSocket和消息相关日志
	Wslog      *zap.SugaredLogger // WebSocket连接日志
	Nsqlog     *zap.SugaredLogger // 消息队列日志
	Collectlog *zap.SugaredLogger // 数据采集日志
	Updatelog  *zap.SugaredLogger // 数据更新日志

	// 业务域专用日志（继承go-novel的设计，便于后续扩展）
	Biquge34log  *zap.SugaredLogger // 业务模块1专用日志
	Paoshu8log   *zap.SugaredLogger // 业务模块2专用日志
	Xswlog       *zap.SugaredLogger // 业务模块3专用日志
	Lydlog       *zap.SugaredLogger // 业务模块4专用日志
	Bqg24log     *zap.SugaredLogger // 业务模块5专用日志
	Siluke520log *zap.SugaredLogger // 业务模块6专用日志
)
