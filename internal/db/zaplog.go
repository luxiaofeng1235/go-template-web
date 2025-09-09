package db

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/utils/zaplog"
)

func InitZapLog() {
	logPath := fmt.Sprintf("%v", "./logs/")
	global.Sqllog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "sql"))
	global.Paylog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "pay"))
	global.Errlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "err"))
	global.Wslog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "ws"))
	global.Collectlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "collect"))
	global.Nsqlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "nsq"))
	global.Updatelog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "update"))
	global.Requestlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "request"))
	global.Biquge34log = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "biquge34"))
	global.Paoshu8log = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "paoshu8"))
	global.Xswlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "xsw"))
	global.Lydlog = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "lyd"))
	global.Bqg24log = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "bqg24"))
	global.Siluke520log = zaplog.LogConfig(fmt.Sprintf("%s%s", logPath, "siluke520"))
	return
}
