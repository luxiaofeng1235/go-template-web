package models

// UserModel 用户数据模型
type UserModel struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
	Email     string `json:"email" gorm:"type:varchar(100);not null"`
	Password  string `json:"-" gorm:"type:varchar(255);not null"` // 密码不返回给前端
	Status    int    `json:"status" gorm:"type:tinyint;default:1;comment:状态:1正常,0禁用"`
	Avatar    string `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "ls_user"
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username string `json:"username" v:"required|length:3,20#用户名不能为空|用户名长度为3-20位"`
	Password string `json:"password" v:"required|length:6,20#密码不能为空|密码长度为6-20位"`
	Email    string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确"`
}

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	Username string `json:"username" v:"required#用户名不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
}

// UserLoginRes 用户登录响应
type UserLoginRes struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Token    string `json:"token"`
}

// UserProfileRes 用户资料响应
type UserProfileRes struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
}
