package server

import (
	"go-web-template/internal/bootstrap"
)

// StartAPIServer 启动API服务器 - 调用bootstrap层的启动逻辑
func StartAPIServer() {
	bootstrap.StartAPIServer()
}

// StartAdminServer 启动管理后台服务器 - 调用bootstrap层的启动逻辑
func StartAdminServer() {
	bootstrap.StartAdminServer()
}
