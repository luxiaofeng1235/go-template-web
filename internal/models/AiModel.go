/*
 * @file: AiModel.go
 * @description: AI服务相关数据模型定义
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package models

// ToImageReq 生成图片请求参数
type ToImageReq struct {
	UserID    int    `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	Model     int    `json:"model" form:"model" v:"required|in:1,2#模型必须选择|模型参数错误"`
	Prompt    string `json:"prompt" form:"prompt" v:"required|max-length:500#提示词必须输入|提示词长度不能超过500字符"`
	Size      string `json:"size" form:"size" v:"required#图片尺寸必须选择"`
	N         int    `json:"n" form:"n" v:"required|between:1,4#生成数量必须输入|生成数量必须在1-4之间"`
	Watermark string `json:"watermark" form:"watermark" v:"required#水印URL必须输入"`
}

// ToVideoReq 生成视频请求参数
type ToVideoReq struct {
	UserID  int    `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	To      int    `json:"to" form:"to" v:"required|in:1,2,3#视频类型必须选择|视频类型参数错误"`
	Prompt  string `json:"prompt" form:"prompt" v:"max-length:500#提示词长度不能超过500字符"`
	ImgURL  string `json:"img_url" form:"img_url"`
}

// GetImageReq 获取图片生成结果请求参数
type GetImageReq struct {
	UserID int    `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	TaskID string `json:"task_id" form:"task_id" v:"required#任务ID必须输入"`
}

// GetVideoReq 获取视频生成结果请求参数
type GetVideoReq struct {
	UserID int    `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	TaskID string `json:"task_id" form:"task_id" v:"required#任务ID必须输入"`
}

// AIWorkListReq AI作品列表请求参数
type AIWorkListReq struct {
	UserID int `json:"user_id" form:"user_id" v:"required#用户ID必须输入"`
	Page   int `json:"page" form:"page" v:"required|min:1#页码必须输入|页码必须大于等于1"`
	Type   int `json:"type" form:"type" v:"required|in:1,2,3#类型必须选择|类型参数错误"`
}

// AIGenerateResult AI生成结果
type AIGenerateResult struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
	URL    string `json:"url,omitempty"`
}