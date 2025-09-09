package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"go-web-template/db"
)

// 管理后台服务器启动入口 - 参考go-novel的架构设计
func main() {
	db.StartAdminServer()
}
