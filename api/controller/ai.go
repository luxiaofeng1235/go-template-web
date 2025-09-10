/*
 * @file: ai.go
 * @description: AI服务API控制器 - 实现图片生成、视频生成等AI功能
 * @author: red
 * @created: 2025-09-10
 * @version: 1.0.0
 * @license: MIT License
 */
package controller

import (
	"fmt"
	"go-web-template/global"
	"go-web-template/internal/models"
	"go-web-template/internal/service/common"
	"go-web-template/utils"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AiController struct{}

// ToImage 生成图片接口
func (c *AiController) ToImage(r *ghttp.Request) {
	var req *models.ToImageReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 验证size格式：应该是 "width~height" 的形式
	sizeParts := strings.Split(req.Size, "~")
	if len(sizeParts) != 2 {
		utils.FailEncrypt(r, fmt.Errorf("size格式错误"), "size格式错误")
		return
	}

	// 检查尺寸是否为数字且在合理范围内
	var width, height int
	var err error
	if width, err = strconv.Atoi(sizeParts[0]); err != nil {
		utils.FailEncrypt(r, err, "size格式错误")
		return
	}
	if height, err = strconv.Atoi(sizeParts[1]); err != nil {
		utils.FailEncrypt(r, err, "size格式错误")
		return
	}

	// 验证尺寸范围：720-1440，总像素不超过200万
	if width < 720 || width > 1440 || height < 720 || height > 1440 {
		utils.FailEncrypt(r, fmt.Errorf("尺寸范围错误"), "尺寸范围必须在720-1440之间")
		return
	}
	if width*height > 2000000 {
		utils.FailEncrypt(r, fmt.Errorf("尺寸过大"), "尺寸过大，像素不能超过200万")
		return
	}

	// 构建图片尺寸字符串
	imageSize := fmt.Sprintf("%d*%d", width, height)
	
	// 调用图片生成服务 - 统一在service层处理复杂逻辑
	result, err := common.GenerateImageByModel(req.Model, req.Prompt, imageSize, req.N, req.Watermark)
	if err != nil {
		global.Errlog.Error(r.Context(), "图片生成失败: %v", err)
		utils.FailEncrypt(r, err, "图片生成失败")
		return
	}

	utils.Success(r, result, "图片生成请求已提交")
}

// GetImage 获取图片生成结果
func (c *AiController) GetImage(r *ghttp.Request) {
	var req *models.GetImageReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 这里应该根据task_id查询任务状态
	// 由于没有具体的任务状态查询URL，先返回模拟数据
	result := &models.AIGenerateResult{
		TaskID: req.TaskID,
		Status: "completed",
		URL:    "https://example.com/generated_image.jpg",
	}

	utils.Success(r, result, "获取图片结果成功")
}

// ToVideo 生成视频接口
func (c *AiController) ToVideo(r *ghttp.Request) {
	var req *models.ToVideoReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 根据to参数进行额外验证
	if req.To == 1 {
		// 图生视频，需要img_url
		if req.ImgURL == "" {
			utils.FailEncrypt(r, fmt.Errorf("img_url必须输入"), "图生视频时img_url必须输入")
			return
		}
	} else {
		// 文生视频，需要prompt
		if req.Prompt == "" {
			utils.FailEncrypt(r, fmt.Errorf("prompt必须输入"), "文生视频时prompt必须输入")
			return
		}
	}

	// 调用视频生成服务 - 统一在service层处理复杂逻辑
	result, err := common.GenerateVideoByType(req.To, req.Prompt, req.ImgURL)
	if err != nil {
		global.Errlog.Error(r.Context(), "视频生成失败: %v", err)
		utils.FailEncrypt(r, err, "视频生成失败")
		return
	}

	utils.Success(r, result, "视频生成请求已提交")
}

// GetVideo 获取视频生成结果
func (c *AiController) GetVideo(r *ghttp.Request) {
	var req *models.GetVideoReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// 这里应该根据task_id查询任务状态
	// 由于没有具体的任务状态查询URL，先返回模拟数据
	result := &models.AIGenerateResult{
		TaskID: req.TaskID,
		Status: "completed",
		URL:    "https://example.com/generated_video.mp4",
	}

	utils.Success(r, result, "获取视频结果成功")
}

// GetAiWorkList 获取AI作品历史记录
func (c *AiController) GetAiWorkList(r *ghttp.Request) {
	var req *models.AIWorkListReq
	if err := r.Parse(&req); err != nil {
		utils.FailEncrypt(r, err, "参数解析失败")
		return
	}

	// 参数验证
	if err := g.Validator().Data(req).Run(r.Context()); err != nil {
		utils.FailEncrypt(r, err, "参数验证失败")
		return
	}

	// TODO: 实现获取AI作品列表的逻辑
	// 这里先返回模拟数据
	result := map[string]interface{}{
		"total": 0,
		"list":  []interface{}{},
		"page":  req.Page,
		"type":  req.Type,
	}

	utils.Success(r, result, "获取AI作品列表成功")
}