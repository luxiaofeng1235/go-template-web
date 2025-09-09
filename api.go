package main

import (
	"go-web-template/db"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

// 测试
// API服务器启动入口 - 参考go-novel的架构设计
func main() {
	db.StartAPIServer()
}
