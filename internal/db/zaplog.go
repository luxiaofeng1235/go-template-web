package db

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/utils/zaplog"
)

func InitZapLog() {
	logPath := fmt.Sprintf("%v", "./logs/")
	global.Sqllog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "sql"))
	global.Errlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "err"))
	global.Wslog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "ws"))
	global.Requestlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "request"))
	return
}
