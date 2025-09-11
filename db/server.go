package db

import (
	"go-web-template/internal/bootstrap"
)

// StartAPIServer 启动API服务器 - 参考go-novel的db包设计模式
func StartAPIServer() {
	bootstrap.StartAPIServer()
}

// StartAdminServer 启动管理后台服务器 - 参考go-novel的db包设计模式
func StartAdminServer() {
	bootstrap.StartAdminServer()
}
