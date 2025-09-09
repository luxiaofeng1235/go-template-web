package utils

import (
	"fmt"
	"strings"
	"github.com/gogf/gf/v2/net/ghttp"
)

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

func GetFullDomain(r *ghttp.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.GetHost()
}
