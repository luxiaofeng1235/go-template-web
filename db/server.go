package db

import (
	"go-web-template/internal/server"
)

// StartAPIServer 启动API服务器 - 参考go-novel的db包设计模式
func StartAPIServer() {
	server.StartAPIServer()
}

// StartAdminServer 启动管理后台服务器 - 参考go-novel的db包设计模式
func StartAdminServer() {
	server.StartAdminServer()
}
