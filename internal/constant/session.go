package constant

/*
 * @Description: 会话配置常量
 * @Version: 1.0.0
 * @Author: red
 * @Date: 2025-09-10
 * @LastEditTime: 2025-09-10
 */

// 用户状态常量 (SecretKey)
const (
	SecretKeyStatusNormal int8 = 1 // 正常
)

// 会话状态常量
const (
	SessionStatusActive  int8 = 1 // 活跃
	SessionStatusOffline int8 = 2 // 离线
	SessionStatusTimeout int8 = 3 // 超时
)
