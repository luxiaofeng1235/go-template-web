package utils

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// 获取对应的静态资源的URL
func ProcessLogoURLForStatic(relativePath string, r *ghttp.Request) string {
	if relativePath == "" {
		return ""
	}
	ctx := gctx.New()
	staticPort := g.Cfg().MustGet(ctx, "server.source.address", ":8082").String()
	domain := GetFullDomain(r)
	return domain + staticPort + relativePath
}

// 获取单独的URL信息
func GetFullDomain(r *ghttp.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.GetHost()
}
