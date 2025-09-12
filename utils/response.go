package utils

import (
	"context"
	"go-web-template/internal/constant"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ResponseData 基础响应序列化器（参考go-novel的Response结构）
type ResponseData struct {
	Code  int         `json:"code"`
	Show  int         `json:"show"`
	Data  interface{} `json:"data,omitempty"`
	Key   string      `json:"key,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// PageResult 分页响应结构（参考go-novel的PageResult）
type PageResult struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// TokenData 带有token的Data结构（参考go-novel的TokenData）
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// Success 成功响应（参考go-novel的Success）
func Success(r *ghttp.Request, data interface{}, msg string) {
	response := ResponseData{
		Code: constant.SUCCESS_CODE,
		Show: 0,
		Data: data,
		Msg:  msg,
	}
	r.Response.WriteJsonExit(response)
}

// SuccessWithShow 带show参数的成功响应（用于视频生成中等特殊状态）
func SuccessWithShow(r *ghttp.Request, data interface{}, msg string) {
	response := ResponseData{
		Code: constant.SUCCESS_CODE,
		Show: 1,
		Data: data,
		Msg:  msg,
	}
	r.Response.WriteJsonExit(response)
}

// SuccessEncrypt 加密成功响应（参考go-novel的SuccessEncrypt）
func SuccessEncrypt(r *ghttp.Request, data interface{}, msg string) {
	response := ResponseData{
		Code: constant.SUCCESS_CODE,
		Show: 0,
		Data: data,
		Msg:  msg,
	}

	// 检查是否需要加密
	isEncrypt := GetApiEncrypt()
	if isEncrypt {
		// AES加密响应数据
		encrypt, err := AesEncryptByCFB(GetApiAesKey(), JSONString(response))
		if err != nil {
			// 加密失败，返回普通响应
			r.Response.WriteJsonExit(response)
			return
		}
		res := g.Map{
			"data": encrypt,
		}
		r.Response.WriteJsonExit(res)
		return
	}

	r.Response.WriteJsonExit(response)
}

// Fail 错误响应（参考go-novel的Fail）
func Fail(r *ghttp.Request, err error, msg string) {
	if err != nil {
		msg = err.Error()
	}
	response := ResponseData{
		Code: constant.ERROR_CODE,
		Show: 0,
		Msg:  msg,
	}
	r.Response.WriteJsonExit(response)
}

// FailEncrypt 加密错误响应（参考go-novel的FailEncrypt）
func FailEncrypt(r *ghttp.Request, err error, msg string) {
	if err != nil {
		msg = err.Error()
	}
	response := ResponseData{
		Code: constant.ERROR_CODE,
		Show: 0,
		Msg:  msg,
	}

	isEncrypt := GetApiEncrypt()
	if isEncrypt {
		encrypt, encErr := AesEncryptByCFB(GetApiAesKey(), JSONString(response))
		if encErr != nil {
			r.Response.WriteJsonExit(response)
			return
		}
		res := g.Map{
			"data": encrypt,
		}
		r.Response.WriteJsonExit(res)
		return
	}

	r.Response.WriteJsonExit(response)
}

// ParamError 参数错误响应
func ParamError(r *ghttp.Request, msg ...string) {
	message := constant.MSG_PARAM_ERROR
	if len(msg) > 0 {
		message = msg[0]
	}

	r.Response.WriteJsonExit(ResponseData{
		Code: constant.ERROR_CODE,
		Show: 0,
		Msg:  message,
		Data: nil,
	})
}

// AuthError 认证错误响应
func AuthError(r *ghttp.Request, msg ...string) {
	message := constant.MSG_AUTH_ERROR
	if len(msg) > 0 {
		message = msg[0]
	}

	r.Response.WriteJsonExit(ResponseData{
		Code: constant.AUTH_ERROR,
		Show: 0,
		Msg:  message,
		Data: nil,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(r *ghttp.Request, list interface{}, total int64, page, pageSize int, msg string) {
	pageData := PageResult{
		Data:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	Success(r, pageData, msg)
}

// GetApiEncrypt 获取API加密配置（从配置文件读取）
func GetApiEncrypt() bool {
	ctx := context.Background()
	return g.Cfg().MustGet(ctx, "api.encrypt", false).Bool()
}

// GetApiAesKey 获取API AES密钥
func GetApiAesKey() string {
	ctx := context.Background()
	return g.Cfg().MustGet(ctx, "api.aesKey", "go-web-template-aes-key").String()
}
