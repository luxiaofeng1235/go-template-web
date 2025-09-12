package utils

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// 获取单独的配置信息
func GetStaticResourceURL(r *ghttp.Request, relativePath string) string {
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	host := r.GetHost()
	if strings.Contains(host, ":") {
		hostParts := strings.Split(host, ":")
		host = hostParts[0]
	}

	return fmt.Sprintf("%s://%s:8082%s", protocol, host, relativePath)
}

// 获取单独的URL信息
func GetFullDomain(r *ghttp.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.GetHost()
}

// 获取对应的静态资源的URL
func ProcessLogoURLForStatic(logo string, r *ghttp.Request) string {
	if logo == "" {
		return ""
	}
	ctx := gctx.New()
	staticPort := g.Cfg().MustGet(ctx, "server.source.address", ":8082").String()
	domain := GetFullDomain(r)
	return domain + staticPort + logo
}
