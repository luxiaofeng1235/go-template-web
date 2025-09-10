/*
 * @file: cors.go
 * @description: CORS跨域处理中间件 - 支持完整的跨域配置
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域中间件
func CORS() func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		// 设置允许的源
		r.Response.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 设置允许的请求头
		r.Response.Header().Set("Access-Control-Allow-Headers", 
			"Authorization, Sec-Fetch-Mode, DNT, X-Mx-ReqToken, Keep-Alive, User-Agent, If-Match, If-None-Match, If-Unmodified-Since, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Accept-Language, Origin, Accept-Encoding, Access-Token, token")
		
		// 设置允许的请求方法
		r.Response.Header().Set("Access-Control-Allow-Methods", 
			"GET, POST, PATCH, PUT, DELETE, OPTIONS")
		
		// 设置预检请求缓存时间（秒）
		r.Response.Header().Set("Access-Control-Max-Age", "1728000")
		
		// 允许携带凭证
		r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 处理OPTIONS预检请求
		if r.Method == "OPTIONS" {
			r.Response.WriteStatusExit(200)
			return
		}
		
		// 继续处理请求
		r.Middleware.Next()
	}
}

// SimpleCORS 简化的CORS中间件
func SimpleCORS() func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		r.Response.Header().Set("Access-Control-Allow-Origin", "*")
		r.Response.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, token")
		r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		if r.Method == "OPTIONS" {
			r.Response.WriteStatusExit(200)
			return
		}
		
		r.Middleware.Next()
	}
}