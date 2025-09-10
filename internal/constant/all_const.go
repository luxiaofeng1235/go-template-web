package constant

// HTTP状态码
const (
	SUCCESS_CODE     = 1   // 成功
	SUCCESS_CODE_200 = 1   //默认200的返回值（预留）
	ERROR_CODE       = 0   // 服务器错误
	AUTH_ERROR       = 401 // 认证失败
	PARAM_ERROR      = 400 // 参数错误
)

// 用户状态
const (
	USER_STATUS_NORMAL  = 1 // 正常
	USER_STATUS_DISABLE = 0 // 禁用
)

// 产品状态
const (
	PRODUCT_STATUS_NORMAL  = 1 // 正常
	PRODUCT_STATUS_DISABLE = 0 // 禁用
)

// JWT相关
const (
	JWT_HEADER_KEY = "Authorization"
	JWT_PREFIX     = "Bearer "
)

// Redis Key前缀
const (
	REDIS_USER_PREFIX    = "user:"
	REDIS_TOKEN_PREFIX   = "token:"
	REDIS_PRODUCT_PREFIX = "product:"
)

// 响应消息
const (
	MSG_SUCCESS        = "操作成功"
	MSG_ERROR          = "操作失败"
	MSG_PARAM_ERROR    = "参数错误"
	MSG_AUTH_ERROR     = "认证失败"
	MSG_TOKEN_INVALID  = "Token无效"
	MSG_USER_NOT_FOUND = "用户不存在"
	MSG_PASSWORD_ERROR = "密码错误"
	MSG_USER_EXIST     = "用户已存在"
)

// 数据库表名
const (
	TABLE_USER    = "ls_user"
	TABLE_PRODUCT = "ls_product"
)
